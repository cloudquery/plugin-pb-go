package managedplugin

import "github.com/rs/zerolog"

type Option func(*Client)

func WithLogger(logger zerolog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithDirectory(directory string) Option {
	return func(c *Client) {
		c.directory = directory
	}
}

func WithNoSentry() Option {
	return func(c *Client) {
		c.noSentry = true
	}
}

func WithNoExec() Option {
	return func(c *Client) {
		c.noExec = true
	}
}

func WithNoProgress() Option {
	return func(c *Client) {
		c.noProgress = true
	}
}

func WithOtelEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.otelEndpoint = endpoint
	}
}

func WithOtelEndpointInsecure() Option {
	return func(c *Client) {
		c.otelEndpointInsecure = true
	}
}

func WithAuthToken(authToken string) Option {
	return func(c *Client) {
		c.authToken = authToken
	}
}

func WithTeamName(teamName string) Option {
	return func(c *Client) {
		c.teamName = teamName
	}
}

func WithLicenseFile(licenseFile string) Option {
	return func(c *Client) {
		c.licenseFile = licenseFile
	}
}

func WithCloudQueryDockerHost(dockerHost string) Option {
	return func(c *Client) {
		c.cqDockerHost = dockerHost
	}
}

func WithUseTCP() Option {
	return func(c *Client) {
		c.useTCP = true
	}
}
