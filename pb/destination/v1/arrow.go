package destination

import (
	"bytes"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
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
