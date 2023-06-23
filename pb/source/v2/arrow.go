package source

import (
	"bytes"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
)

const (
	MetadataTableName        = "cq:table_name"
)

type Schemas []*arrow.Schema

func (s Schemas) Len() int {
	return len(s)
}

func (s Schemas) SchemaByName(name string) *arrow.Schema {
	for _, sc := range s {
		tableName, ok := sc.Metadata().GetValue(MetadataTableName)
		if !ok {
			continue
		}
		if tableName == name {
			return sc
		}
	}
	return nil
}

func (s Schemas) Encode() ([][]byte, error) {
	ret := make([][]byte, len(s))
	for i, sc := range s {
		var buf bytes.Buffer
		wr := ipc.NewWriter(&buf, ipc.WithSchema(sc))
		if err := wr.Close(); err != nil {
			return nil, err
		}
		ret[i] = buf.Bytes()
	}
	return ret, nil
}