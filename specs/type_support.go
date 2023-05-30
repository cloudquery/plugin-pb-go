package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TypeSupport int

const (
	TypeSupportLimited TypeSupport = iota
	TypeSupportFull
)

func (r TypeSupport) String() string {
	return [...]string{"limited", "full"}[r]
}
func (r TypeSupport) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (r *TypeSupport) UnmarshalJSON(data []byte) (err error) {
	var typeSupport string
	if err := json.Unmarshal(data, &typeSupport); err != nil {
		return err
	}
	if *r, err = TypeSupportFromString(typeSupport); err != nil {
		return err
	}
	return nil
}

func TypeSupportFromString(s string) (TypeSupport, error) {
	switch s {
	case "limited":
		return TypeSupportLimited, nil
	case "full":
		return TypeSupportFull, nil
	default:
		return TypeSupportLimited, fmt.Errorf("unknown type support %s", s)
	}
}
