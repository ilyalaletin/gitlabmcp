package groups

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

func TestGetGroup_MissingGroupID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getGroup(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing group_id")
	}
}

func TestListGroupProjects_MissingGroupID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listGroupProjects(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing group_id")
	}
}

func TestListGroupMembers_MissingGroupID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listGroupMembers(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing group_id")
	}
}
