package env

import "os"

func TenantID() string {
	return os.Getenv("CQ_TENANT_ID")
}
