package managedplugin

import (
	"context"
	"net"
	"net/http"
	"time"
)

var (
	downloadDialTimeout           = 10 * time.Second
	downloadTLSHandshakeTimeout   = 10 * time.Second
	downloadResponseHeaderTimeout = 30 * time.Second
	downloadIdleReadTimeout       = 30 * time.Second
)

var downloadClient = newDownloadClient()

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
