package issues

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

func TestGetIssue_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getIssue(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetIssue_MissingIID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getIssue(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing issue_iid")
	}
}

func TestCreateIssue_MissingTitle(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createIssue(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing title")
	}
}
