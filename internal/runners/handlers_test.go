package runners

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

func TestGetRunner_MissingRunnerID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getRunner(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing runner_id")
	}
}

func TestEnableRunner_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.enableRunner(context.Background(), makeReq(map[string]any{"runner_id": 123}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestEnableRunner_MissingRunnerID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.enableRunner(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing runner_id")
	}
}

func TestDisableRunner_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.disableRunner(context.Background(), makeReq(map[string]any{"runner_id": 123}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestDisableRunner_MissingRunnerID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.disableRunner(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing runner_id")
	}
}

func TestDeleteRunner_MissingRunnerID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.deleteRunner(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing runner_id")
	}
}
