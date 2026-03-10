package pipelines

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/ilya/gitlabmcp/internal/handler"
	"github.com/mark3labs/mcp-go/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type handlers struct {
	gl *gitlab.Client
}

// All handlers use mcp-go's built-in typed getters on CallToolRequest:
//   req.RequireString("key") -> (string, error)
//   req.GetString("key", "default") -> string
//   req.RequireInt("key") -> (int64, error)
//   req.GetInt("key", defaultVal) -> int64
//   req.GetBool("key", defaultVal) -> bool
// GitLab client-go uses int64 for all IDs and pagination fields.

func (h *handlers) listPipelines(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if s := req.GetString("status", ""); s != "" {
		opts.Status = gitlab.Ptr(gitlab.BuildStateValue(s))
	}
	if s := req.GetString("ref", ""); s != "" {
		opts.Ref = gitlab.Ptr(s)
	}

	pipelines, resp, err := h.gl.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		slog.Error("list_pipelines failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_pipelines", "project", projectID, "count", len(pipelines), "duration", time.Since(start))
	result := handler.NewPaginatedResult(pipelines, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getPipeline(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	pipelineID, err := req.RequireInt("pipeline_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pipeline, _, err := h.gl.Pipelines.GetPipeline(projectID, int64(pipelineID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(pipeline)), nil
}

func (h *handlers) listPipelineJobs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	pipelineID, err := req.RequireInt("pipeline_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	jobs, resp, err := h.gl.Jobs.ListPipelineJobs(projectID, int64(pipelineID), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	result := handler.NewPaginatedResult(jobs, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getJobLog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	jobID, err := req.RequireInt("job_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	reader, _, err := h.gl.Jobs.GetTraceFile(projectID, int64(jobID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	// Read all bytes from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	// Truncate if larger than 50KB
	const maxSize = 50 * 1024
	if len(data) > maxSize {
		data = data[:maxSize]
	}

	return mcp.NewToolResultText(string(data)), nil
}

func (h *handlers) retryPipeline(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	pipelineID, err := req.RequireInt("pipeline_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pipeline, _, err := h.gl.Pipelines.RetryPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(pipeline)), nil
}

func (h *handlers) cancelPipeline(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	pipelineID, err := req.RequireInt("pipeline_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pipeline, _, err := h.gl.Pipelines.CancelPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(pipeline)), nil
}

func (h *handlers) listCIVariables(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectVariablesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	variables, resp, err := h.gl.ProjectVariables.ListVariables(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	result := handler.NewPaginatedResult(variables, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) createCIVariable(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	key, err := req.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	value, err := req.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.CreateProjectVariableOptions{
		Key:   gitlab.Ptr(key),
		Value: gitlab.Ptr(value),
	}

	protected := req.GetBool("protected", false)
	if protected {
		opts.Protected = gitlab.Ptr(true)
	}

	masked := req.GetBool("masked", false)
	if masked {
		opts.Masked = gitlab.Ptr(true)
	}

	if s := req.GetString("environment_scope", ""); s != "" {
		opts.EnvironmentScope = gitlab.Ptr(s)
	}

	variable, _, err := h.gl.ProjectVariables.CreateVariable(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(variable)), nil
}

func (h *handlers) updateCIVariable(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	key, err := req.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.UpdateProjectVariableOptions{}

	if s := req.GetString("value", ""); s != "" {
		opts.Value = gitlab.Ptr(s)
	}

	if p := req.GetBool("protected", false); p {
		opts.Protected = gitlab.Ptr(true)
	}

	if m := req.GetBool("masked", false); m {
		opts.Masked = gitlab.Ptr(true)
	}

	variable, _, err := h.gl.ProjectVariables.UpdateVariable(projectID, key, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(variable)), nil
}

func (h *handlers) deleteCIVariable(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	key, err := req.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err = h.gl.ProjectVariables.RemoveVariable(projectID, key, nil)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(`{"status": "deleted"}`), nil
}
