package specs

import (
	"fmt"
)

type Registry int

const (
	RegistryGithub Registry = iota
	RegistryLocal
	RegistryGrpc
)

func (r Registry) String() string {
	return [...]string{"github", "local", "grpc"}[r]
}

func RegistryFromString(s string) (Registry, error) {
	switch s {
	case "github":
		return RegistryGithub, nil
	case "local":
		return RegistryLocal, nil
	case "grpc":
		return RegistryGrpc, nil
	default:
		return RegistryGithub, fmt.Errorf("unknown registry %s", s)
	}
}
