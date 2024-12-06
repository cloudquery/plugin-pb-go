package source

import (
	"bytes"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
)

func SchemasToBytes(schemas []*arrow.Schema) ([][]byte, error) {
	ret := make([][]byte, len(schemas))
	for i, sc := range schemas {
		var buf bytes.Buffer
		wr := ipc.NewWriter(&buf, ipc.WithSchema(sc))
		if err := wr.Close(); err != nil {
			return nil, err
		}
		ret[i] = buf.Bytes()
	}
	return ret, nil
}

func RecordToBytes(record arrow.Record) ([]byte, error) {
	var buf bytes.Buffer
	wr := ipc.NewWriter(&buf, ipc.WithSchema(record.Schema()))
	if err := wr.Write(record); err != nil {
		return nil, err
	}
	if err := wr.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
