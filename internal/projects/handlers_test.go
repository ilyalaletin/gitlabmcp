package projects

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

func TestGetProject_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getProject(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestCreateProject_MissingName(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createProject(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing name")
	}
}

func TestListProjectMembers_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listProjectMembers(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestListProjectWebhooks_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listProjectWebhooks(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestCreateProjectWebhook_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createProjectWebhook(context.Background(), makeReq(map[string]any{"url": "https://example.com/webhook"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestCreateProjectWebhook_MissingURL(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createProjectWebhook(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing url")
	}
}
