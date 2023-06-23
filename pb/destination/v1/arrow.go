package destination

import (
	"bytes"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
)

const (
	MetadataTableName = "cq:table_name"
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

func NewSchemasFromBytes(b [][]byte) (Schemas, error) {
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
