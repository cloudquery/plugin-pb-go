package managedplugin

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type DownloadSource int

const (
	DownloadSourceUnknown DownloadSource = iota
	DownloadSourceCached
	DownloadSourceRemote
)

func (r DownloadSource) String() string {
	return [...]string{"unknown", "cached", "remote"}[r]
}

func (r DownloadSource) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, r.String())), nil
}

func (r *DownloadSource) UnmarshalJSON(data []byte) (err error) {
	var mode string
	if err := json.Unmarshal(data, &mode); err != nil {
		return err
	}
	if *r, err = DownloadModeFromString(mode); err != nil {
		return err
	}
	return nil
}

func DownloadModeFromString(s string) (DownloadSource, error) {
	switch s {
	case "cached":
		return DownloadSourceCached, nil
	case "remote":
		return DownloadSourceRemote, nil
	default:
		return DownloadSourceUnknown, fmt.Errorf("unknown registry %s", s)
	}
}

type Metrics struct {
	Errors         uint64
	Warnings       uint64
	DownloadSource DownloadSource
}

func (m *Metrics) incrementErrors() {
	atomic.AddUint64(&m.Errors, 1)
}

func (m *Metrics) incrementWarnings() {
	atomic.AddUint64(&m.Warnings, 1)
}
