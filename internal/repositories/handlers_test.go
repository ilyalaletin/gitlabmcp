package repositories

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

func TestListBranches_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listBranches(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetFile_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getFile(context.Background(), makeReq(map[string]any{
		"file_path": "README.md",
		"ref":       "main",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetFile_MissingFilePath(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getFile(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
		"ref":        "main",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing file_path")
	}
}

func TestGetFile_MissingRef(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getFile(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
		"file_path":  "README.md",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing ref")
	}
}

func TestListRepositoryTree_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listRepositoryTree(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestListCommits_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listCommits(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetCommit_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getCommit(context.Background(), makeReq(map[string]any{
		"sha": "abc123",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetCommit_MissingSHA(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getCommit(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing sha")
	}
}
