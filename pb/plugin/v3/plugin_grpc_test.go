package plugin //nolint:revive // var-naming: package name conflicts with standard library but is appropriate for this context

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnimplementedAssessTables(t *testing.T) {
	_, err := UnimplementedPluginServer{}.AssessTables(context.Background(), &AssessTables_Request{})
	if status.Code(err) != codes.Unimplemented {
		t.Errorf("expected %v, got %v", codes.Unimplemented, err)
	}
}
