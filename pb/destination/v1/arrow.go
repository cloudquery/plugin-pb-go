package destination

import (
	"bytes"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/ipc"
)

func NewSchemasFromBytes(b [][]byte) ([]*arrow.Schema, error) {
	ret := make([]*arrow.Schema, len(b))
	for i, buf := range b {
		rdr, err := ipc.NewReader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		ret[i] = rdr.Schema()
	}
	return ret, nil
}

func NewRecordFromBytes(b []byte) (arrow.Record, error) {
	rdr, err := ipc.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	for rdr.Next() {
		rec := rdr.Record()
		rec.Retain()
		return rec, nil
	}
	return nil, nil
}
