package managedplugin

import "github.com/rs/zerolog"

func (c *Client) jsonToLog(l zerolog.Logger, msg map[string]any) {
	level := msg["level"]
	// The log level is part of the log message received from the plugin, so we need to remove it before logging
	delete(msg, "level")
	switch level {
	case "trace":
		l.Trace().Fields(msg).Msg("")
	case "debug":
		l.Debug().Fields(msg).Msg("")
	case "info":
		l.Info().Fields(msg).Msg("")
	case "warn":
		l.Warn().Fields(msg).Msg("")
		c.metrics.incrementWarnings()
	case "error":
		l.Error().Fields(msg).Msg("")
		if msg["message"].(string) == "table resolver finished with error" && msg["table"].(string) != "" {
			c.metrics.addErroredTable(msg["table"].(string))
		}
		c.metrics.incrementErrors()
	default:
		l.Error().Fields(msg).Msg("unknown level")
	}
}
