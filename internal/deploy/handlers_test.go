package deploy

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// Tests verify parameter validation. Full API tests require integration setup.
// Helper to build a CallToolRequest with given arguments.
func makeReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestGetEnvironment_MissingEnvironmentID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getEnvironment(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing environment_id")
	}
}

func TestCreateEnvironment_MissingName(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createEnvironment(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing name")
	}
}

func TestStopEnvironment_MissingEnvironmentID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.stopEnvironment(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing environment_id")
	}
}

func TestListDeployments_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listDeployments(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}
