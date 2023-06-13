package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

// Plugin is the spec for a CloudQuery plugin
type Plugin struct {
	// Name of the source plugin to use
	Name string `json:"name,omitempty"`
	// Version of the source plugin to use
	Version string `json:"version,omitempty"`
	// Path is the canonical path to the source plugin in a given registry
	// For example:
	// in github the path will be: org/repo
	// For the local registry the path will be the path to the binary: ./path/to/binary
	// For the gRPC registry the path will be the address of the gRPC server: host:port
	Path string
	// Registry can be github,local,grpc.
	Registry Registry `json:"registry,omitempty"`
	// Sync defines sync behaviour
	Sync Sync `json:"sync,omitempty"`
	// StateBackend specifies the state backend to use
	// for incremental syncs
	StateBackend *Plugin `json:"backend,omitempty"`
	// Write specifies write behaviour
	Write Write `json:"write,omitempty"`
	// Spec is defined by each plugin
	Spec any `json:"spec,omitempty"`
}

func (d *Plugin) UnmarshalSpec(out any) error {
	b, err := json.Marshal(d.Spec)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func (d *Plugin) VersionString() string {
	if d.Registry != RegistryGithub {
		return fmt.Sprintf("%s (%s@%s)", d.Name, d.Registry, d.Path)
	}
	pathParts := strings.Split(d.Path, "/")
	if len(pathParts) != 2 {
		return fmt.Sprintf("%s (%s@%s)", d.Name, d.Path, d.Version)
	}
	if d.Name == pathParts[1] {
		return fmt.Sprintf("%s (%s)", d.Name, d.Version)
	}
	return fmt.Sprintf("%s (%s@%s)", d.Name, pathParts[1], d.Version)
}

func (d *Plugin) SetDefaults(defaultBatchSize, defaultBatchSizeBytes int) {
	if d.Write.BatchSize == 0 {
		d.Write.BatchSize = defaultBatchSize
	}
	if d.Write.BatchSizeBytes == 0 {
		d.Write.BatchSizeBytes = defaultBatchSizeBytes
	}
}

func (d *Plugin) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.Path == "" {
		msg := "path is required"
		// give a small hint to help users transition from the old config format that didn't require path
		officialPlugins := []string{"postgresql", "csv"}
		if funk.ContainsString(officialPlugins, d.Name) {
			msg += fmt.Sprintf(". Hint: try setting path to cloudquery/%s in your config", d.Name)
		}
		return fmt.Errorf(msg)
	}

	if d.Registry == RegistryGithub {
		if d.Version == "" {
			return fmt.Errorf("version is required")
		}
		if !strings.HasPrefix(d.Version, "v") {
			return fmt.Errorf("version must start with v")
		}
	}
	if d.Write.BatchSize < 0 {
		return fmt.Errorf("batch_size must be greater than 0")
	}
	return nil
}
