package mr

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

func TestGetMergeRequest_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getMergeRequest(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetMergeRequest_MissingIID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getMergeRequest(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing merge_request_iid")
	}
}

func TestCreateMergeRequest_MissingTitle(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createMergeRequest(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing title")
	}
}

func TestCreateMergeRequest_MissingSourceBranch(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createMergeRequest(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
		"title":      "Test MR",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing source_branch")
	}
}

func TestCreateMergeRequest_MissingTargetBranch(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createMergeRequest(context.Background(), makeReq(map[string]any{
		"project_id":    "owner/repo",
		"title":         "Test MR",
		"source_branch": "feature",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing target_branch")
	}
}
