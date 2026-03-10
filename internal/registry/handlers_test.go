package registry

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

func TestListRegistryRepos_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listRegistryRepos(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestListRegistryTags_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listRegistryTags(context.Background(), makeReq(map[string]any{"repository_id": 123}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestListRegistryTags_MissingRepositoryID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listRegistryTags(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing repository_id")
	}
}

func TestDeleteRegistryTag_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.deleteRegistryTag(context.Background(), makeReq(map[string]any{"repository_id": 123, "tag_name": "v1.0.0"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestDeleteRegistryTag_MissingRepositoryID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.deleteRegistryTag(context.Background(), makeReq(map[string]any{"project_id": "owner/repo", "tag_name": "v1.0.0"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing repository_id")
	}
}

func TestDeleteRegistryTag_MissingTagName(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.deleteRegistryTag(context.Background(), makeReq(map[string]any{"project_id": "owner/repo", "repository_id": 123}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing tag_name")
	}
}
