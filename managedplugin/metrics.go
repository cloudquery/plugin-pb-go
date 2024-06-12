package managedplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type DownloadMode int

const (
	DownloadModeUnknown DownloadMode = iota
	DownloadModeCached
	DownloadModeRemote
)

func (r DownloadMode) String() string {
	return [...]string{"unknown", "cached", "remote"}[r]
}

func (r DownloadMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (r *DownloadMode) UnmarshalJSON(data []byte) (err error) {
	var mode string
	if err := json.Unmarshal(data, &mode); err != nil {
		return err
	}
	if *r, err = DownloadModeFromString(mode); err != nil {
		return err
	}
	return nil
}

func DownloadModeFromString(s string) (DownloadMode, error) {
	switch s {
	case "cached":
		return DownloadModeCached, nil
	case "remote":
		return DownloadModeRemote, nil
	default:
		return DownloadModeUnknown, fmt.Errorf("unknown registry %s", s)
	}
}

type Metrics struct {
	Errors       uint64
	Warnings     uint64
	DownloadMode DownloadMode
}

func (m *Metrics) incrementErrors() {
	atomic.AddUint64(&m.Errors, 1)
}

func (m *Metrics) incrementWarnings() {
	atomic.AddUint64(&m.Warnings, 1)
}
