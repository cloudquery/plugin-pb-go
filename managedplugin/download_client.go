package managedplugin

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Overridable so tests do not pay the real timeouts.
var (
	downloadDialTimeout           = 10 * time.Second
	downloadTLSHandshakeTimeout   = 10 * time.Second
	downloadResponseHeaderTimeout = 30 * time.Second
	downloadIdleReadTimeout       = 30 * time.Second
)

var downloadClient = newDownloadClient()

// idleTimeoutConn arms a read deadline immediately before every read, so a peer
// that stops sending fails the connection instead of holding it until some
// middlebox tears it down minutes later. Arming per read rather than once means
// a slow but moving transfer never trips, and a caller stalled in its own write
// path is not blamed on the server.
type idleTimeoutConn struct {
	net.Conn
	idle time.Duration
}

func (c *idleTimeoutConn) Read(b []byte) (int, error) {
	if err := c.SetReadDeadline(time.Now().Add(c.idle)); err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

// newDownloadClient builds the client used for every plugin asset request.
// http.DefaultClient has no timeout of any kind, so a silent server hangs a
// download until the connection is torn down externally.
//
// ForceAttemptHTTP2 is deliberately left off: a read deadline is per connection,
// and HTTP/2 multiplexes streams onto one connection, so an idle deadline there
// would be shared across concurrent requests.
func newDownloadClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   downloadDialTimeout,
		KeepAlive: 30 * time.Second,
	}
	idle := downloadIdleReadTimeout
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return &idleTimeoutConn{Conn: conn, idle: idle}, nil
			},
			TLSHandshakeTimeout:   downloadTLSHandshakeTimeout,
			ResponseHeaderTimeout: downloadResponseHeaderTimeout,
			IdleConnTimeout:       downloadIdleReadTimeout,
			MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
		},
	}
}
