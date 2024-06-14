package managedplugin

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type AssetSource int

const (
	AssetSourceUnknown AssetSource = iota
	AssetSourceCached
	AssetSourceRemote
)

func (r AssetSource) String() string {
	return [...]string{"unknown", "cached", "remote"}[r]
}

func (r AssetSource) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, r.String())), nil
}

func (r *AssetSource) UnmarshalJSON(data []byte) (err error) {
	var mode string
	if err := json.Unmarshal(data, &mode); err != nil {
		return err
	}
	if *r, err = AssetSourceFromString(mode); err != nil {
		return err
	}
	return nil
}

func AssetSourceFromString(s string) (AssetSource, error) {
	switch s {
	case "cached":
		return AssetSourceCached, nil
	case "remote":
		return AssetSourceRemote, nil
	default:
		return AssetSourceUnknown, fmt.Errorf("unknown mode %s", s)
	}
}

type Metrics struct {
	Errors      uint64
	Warnings    uint64
	AssetSource AssetSource
}

func (m *Metrics) incrementErrors() {
	atomic.AddUint64(&m.Errors, 1)
}

func (m *Metrics) incrementWarnings() {
	atomic.AddUint64(&m.Warnings, 1)
}
