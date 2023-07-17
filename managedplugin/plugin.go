package managedplugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	pbBase "github.com/cloudquery/plugin-pb-go/pb/base/v0"
	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v0"
	pbDiscoveryV1 "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pbSource "github.com/cloudquery/plugin-pb-go/pb/source/v0"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultDownloadDir = ".cq"
	maxMsgSize         = 100 * 1024 * 1024 // 100 MiB
)

// PluginType specifies if a plugin is a source or a destination
// it actually doesn't really have any effect as plugins can serve both as source and as destinations
// but it is here for backward compatibility
type PluginType int

const (
	PluginSource PluginType = iota
	PluginDestination
)

func (p PluginType) String() string {
	return [...]string{"source", "destination"}[p]
}

type Clients []*Client

type Config struct {
	Name     string
	Registry Registry
	Path     string
	Version  string
}

type Client struct {
	directory            string
	cmd                  *exec.Cmd
	logger               zerolog.Logger
	LocalPath            string
	grpcSocketName       string
	wg                   *sync.WaitGroup
	Conn                 *grpc.ClientConn
	config               Config
	noSentry             bool
	otelEndpoint         string
	otelEndpointInsecure bool
	metrics              *Metrics
	registry             Registry
}

// typ will be deprecated soon but now required for a transition period
func NewClients(ctx context.Context, typ PluginType, specs []Config, opts ...Option) (Clients, error) {
	clients := make(Clients, len(specs))
	for i, spec := range specs {
		client, err := NewClient(ctx, typ, spec, opts...)
		if err != nil {
			return nil, err
		}
		clients[i] = client
	}
	return clients, nil
}

func (c Clients) ClientByName(name string) *Client {
	for _, client := range c {
		if client.config.Name == name {
			return client
		}
	}
	return nil
}

func (c Clients) Terminate() error {
	for _, client := range c {
		if err := client.Terminate(); err != nil {
			return err
		}
	}
	return nil
}

// NewClient creates a new plugin client.
// If registrySpec is GitHub then client downloads the plugin, spawns it and creates a gRPC connection.
// If registrySpec is Local then client spawns the plugin and creates a gRPC connection.
// If registrySpec is gRPC then clients creates a new connection
func NewClient(ctx context.Context, typ PluginType, config Config, opts ...Option) (*Client, error) {
	c := Client{
		directory: defaultDownloadDir,
		wg:        &sync.WaitGroup{},
		config:    config,
		metrics:   &Metrics{},
		registry:  config.Registry,
	}
	for _, opt := range opts {
		opt(&c)
	}
	var err error
	switch config.Registry {
	case RegistryGrpc:
		c.Conn, err = grpc.DialContext(ctx, config.Path,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(maxMsgSize),
				grpc.MaxCallSendMsgSize(maxMsgSize),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to dial grpc source plugin at %s: %w", config.Path, err)
		}
	case RegistryLocal:
		if err := validateLocalExecPath(config.Path); err != nil {
			return nil, err
		}
		if err := c.startLocal(ctx, config.Path); err != nil {
			return nil, err
		}
	case RegistryGithub:
		pathSplit := strings.Split(config.Path, "/")
		if len(pathSplit) != 2 {
			return nil, fmt.Errorf("invalid github plugin path: %s. format should be owner/repo", config.Path)
		}
		org, name := pathSplit[0], pathSplit[1]
		c.LocalPath = filepath.Join(c.directory, "plugins", typ.String(), org, name, config.Version, "plugin")
		c.LocalPath = WithBinarySuffix(c.LocalPath)
		if err := DownloadPluginFromGithub(ctx, c.LocalPath, org, name, config.Version, typ); err != nil {
			return nil, err
		}
		if err := c.startLocal(ctx, c.LocalPath); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *Client) ConnectionString() string {
	tgt := c.Conn.Target()
	switch c.registry {
	case RegistryGrpc:
		return tgt
	case RegistryLocal:
		return "unix://" + tgt
	case RegistryGithub:
		return "unix://" + tgt
	}
	return tgt
}

func (c *Client) Metrics() Metrics {
	return *c.metrics
}

func (c *Client) startLocal(ctx context.Context, path string) error {
	c.grpcSocketName = GenerateRandomUnixSocketName()
	// spawn the plugin first and then connect
	args := []string{"serve", "--network", "unix", "--address", c.grpcSocketName,
		"--log-level", c.logger.GetLevel().String(), "--log-format", "json"}
	if c.noSentry {
		args = append(args, "--no-sentry")
	}
	if c.otelEndpoint != "" {
		args = append(args, "--otel-endpoint", c.otelEndpoint)
	}
	if c.otelEndpointInsecure {
		args = append(args, "--otel-endpoint-insecure")
	}
	cmd := exec.CommandContext(ctx, path, args...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = getSysProcAttr()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin %s: %w", path, err)
	}

	c.cmd = cmd

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		lr := NewLogReader(reader)
		for {
			line, err := lr.NextLine()
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, ErrLogLineToLong) {
				c.logger.Info().Str("line", string(line)).Msg("truncated plugin log line")
				continue
			}
			if err != nil {
				c.logger.Err(err).Msg("failed to read log line from plugin")
				break
			}
			var structuredLogLine map[string]any
			if err := json.Unmarshal(line, &structuredLogLine); err != nil {
				c.logger.Err(err).Str("line", string(line)).Msg("failed to unmarshal log line from plugin")
			} else {
				c.jsonToLog(c.logger, structuredLogLine)
			}
		}
	}()

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, "unix", addr)
	}
	c.Conn, err = grpc.DialContext(ctx, c.grpcSocketName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		))
	if err != nil {
		if err := cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill plugin process")
		}
		return err
	}
	return nil
}

func (c *Client) Name() string {
	return c.config.Name
}

func (c *Client) oldDiscovery(ctx context.Context) ([]int, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(c.Conn)
	versionsRes, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		if isUnimplemented(err) {
			// If we get an error here, we assume that the plugin is not a v1 plugin and we try to sync it as a v0 plugin
			// this is for backward compatibility where we used incorrect versioning mechanism
			oldDiscoveryClient := pbSource.NewSourceClient(c.Conn)
			versionRes, err := oldDiscoveryClient.GetProtocolVersion(ctx, &pbBase.GetProtocolVersion_Request{})
			if err != nil {
				return nil, err
			}
			switch versionRes.Version {
			case 2:
				return []int{0}, nil
			case 1:
				return []int{-1}, nil
			default:
				return nil, fmt.Errorf("unknown protocol version %d", versionRes.Version)
			}
		}
		return nil, err
	}
	versions := make([]int, len(versionsRes.Versions))
	for i, vStr := range versionsRes.Versions {
		vStr = strings.TrimPrefix(vStr, "v")
		v, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %w", vStr, err)
		}
		versions[i] = int(v)
	}
	return versions, nil
}

func (c *Client) Versions(ctx context.Context) ([]int, error) {
	discoveryClient := pbDiscoveryV1.NewDiscoveryClient(c.Conn)
	versionsRes, err := discoveryClient.GetVersions(ctx, &pbDiscoveryV1.GetVersions_Request{})
	if err != nil {
		if isUnimplemented(err) {
			// this was only added post v3 so clients will fallback to using an older discovery service
			return c.oldDiscovery(ctx)
		}
		return nil, err
	}
	res := make([]int, len(versionsRes.Versions))
	for i, v := range versionsRes.Versions {
		res[i] = int(v)
	}
	return res, nil
}

func (c *Client) MaxVersion(ctx context.Context) (int, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(c.Conn)
	versionsRes, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		// If we get an error here, we assume that the plugin is not a v1 plugin and we try to sync it as a v0 plugin
		// this is for backward compatibility where we used incorrect versioning mechanism
		oldDiscoveryClient := pbSource.NewSourceClient(c.Conn)
		versionRes, err := oldDiscoveryClient.GetProtocolVersion(ctx, &pbBase.GetProtocolVersion_Request{})
		if err != nil {
			return -1, err
		}
		switch versionRes.Version {
		case 2:
			return 0, nil
		case 1:
			return -1, nil
		default:
			return -1, fmt.Errorf("unknown protocol version %d", versionRes.Version)
		}
	}
	if slices.Contains(versionsRes.Versions, "v2") {
		return 2, nil
	}
	if slices.Contains(versionsRes.Versions, "v1") {
		return 1, nil
	}
	if slices.Contains(versionsRes.Versions, "v0") {
		return 0, nil
	}
	return -1, fmt.Errorf("unknown protocol versions %v", versionsRes.Versions)
}

func (c *Client) Terminate() error {
	// wait for log streaming to complete before returning from this function
	defer c.wg.Wait()

	if c.grpcSocketName != "" {
		defer func() {
			if err := os.RemoveAll(c.grpcSocketName); err != nil {
				c.logger.Error().Err(err).Msg("failed to remove source socket file")
			}
		}()
	}

	if c.Conn != nil {
		if err := c.Conn.Close(); err != nil {
			c.logger.Error().Err(err).Msg("failed to close gRPC connection to source plugin")
		}
		c.Conn = nil
	}
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.terminateProcess(); err != nil {
			return err
		}
	}

	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func validateLocalExecPath(filePath string) error {
	directory, err := isDirectory(filePath)
	if err != nil {
		return fmt.Errorf("error validating plugin path, %s: %w", filePath, err)
	}
	if directory {
		return fmt.Errorf("invalid plugin path: %s. Path cannot point to a directory", filePath)
	}
	return nil
}
