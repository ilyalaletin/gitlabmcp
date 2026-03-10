package pipelines

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

func TestGetPipeline_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getPipeline(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetPipeline_MissingPipelineID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getPipeline(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing pipeline_id")
	}
}

func TestListPipelineJobs_MissingPipelineID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.listPipelineJobs(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing pipeline_id")
	}
}

func TestGetJobLog_MissingJobID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getJobLog(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing job_id")
	}
}

func TestCreateCIVariable_MissingKey(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createCIVariable(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
		"value":      "test_value",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing key")
	}
}

func TestCreateCIVariable_MissingValue(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createCIVariable(context.Background(), makeReq(map[string]any{
		"project_id": "owner/repo",
		"key":        "TEST_KEY",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing value")
	}
}

func TestDeleteCIVariable_MissingKey(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.deleteCIVariable(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing key")
	}
}
