package plugin

import (
	"bytes"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/ipc"
)

func SchemaToBytes(sc *arrow.Schema) ([]byte, error) {
	var buf bytes.Buffer
	wr := ipc.NewWriter(&buf, ipc.WithSchema(sc))
	if err := wr.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewSchemaFromBytes(b []byte) (*arrow.Schema, error) {
	rdr, err := ipc.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return rdr.Schema(), nil
}

func SchemasToBytes(schemas []*arrow.Schema) ([][]byte, error) {
	ret := make([][]byte, len(schemas))
	for i, sc := range schemas {
		buf, err := SchemaToBytes(sc)
		if err != nil {
			return nil, err
		}
		ret[i] = buf
	}
	return ret, nil
}

func NewSchemasFromBytes(b [][]byte) ([]*arrow.Schema, error) {
	schemas := make([]*arrow.Schema, len(b))
	for i, buf := range b {
		sc, err := NewSchemaFromBytes(buf)
		if err != nil {
			return nil, err
		}
		schemas[i] = sc
	}
	return schemas, nil
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
