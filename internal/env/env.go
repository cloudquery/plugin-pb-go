package env

import "os"

func TenantID() string {
	return os.Getenv("_CQ_TENANT_ID")
}
