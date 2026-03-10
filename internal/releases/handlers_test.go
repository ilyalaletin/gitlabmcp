package releases

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

func TestListReleases_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listReleases(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetRelease_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getRelease(context.Background(), makeReq(map[string]any{"tag_name": "v1.0.0"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetRelease_MissingTagName(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getRelease(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing tag_name")
	}
}

func TestCreateRelease_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createRelease(context.Background(), makeReq(map[string]any{"tag_name": "v1.0.0"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestCreateRelease_MissingTagName(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createRelease(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing tag_name")
	}
}
