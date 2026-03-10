package runners

import (
	"context"
	"log/slog"
	"time"

	"github.com/ilyalaletin/gitlabmcp/internal/handler"
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

func (h *handlers) listRunners(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID := req.GetString("project_id", "")
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	var runners []*gitlab.Runner
	var resp *gitlab.Response
	var err error

	if projectID == "" {
		// List all accessible runners
		opts := &gitlab.ListRunnersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    int64(page),
				PerPage: perPageClamped,
			},
		}

		if s := req.GetString("scope", ""); s != "" {
			opts.Scope = gitlab.Ptr(s)
		}

		runners, resp, err = h.gl.Runners.ListRunners(opts)
		if err != nil {
			slog.Error("list_runners failed", "error", err, "duration", time.Since(start))
			return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
		}
		slog.Info("list_runners (global)", "count", len(runners), "duration", time.Since(start))
	} else {
		// List runners for specific project
		opts := &gitlab.ListProjectRunnersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    int64(page),
				PerPage: perPageClamped,
			},
		}

		if s := req.GetString("scope", ""); s != "" {
			opts.Scope = gitlab.Ptr(s)
		}

		runners, resp, err = h.gl.Runners.ListProjectRunners(projectID, opts)
		if err != nil {
			slog.Error("list_runners failed", "project", projectID, "error", err, "duration", time.Since(start))
			return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
		}
		slog.Info("list_runners (project)", "project", projectID, "count", len(runners), "duration", time.Since(start))
	}

	result := handler.NewPaginatedResult(runners, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getRunner(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runnerID, err := req.RequireInt("runner_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	runner, _, err := h.gl.Runners.GetRunnerDetails(int(runnerID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(runner)), nil
}

func (h *handlers) enableRunner(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	runnerID, err := req.RequireInt("runner_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.EnableProjectRunnerOptions{
		RunnerID: int64(runnerID),
	}

	runner, _, err := h.gl.Runners.EnableProjectRunner(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(runner)), nil
}

func (h *handlers) disableRunner(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	runnerID, err := req.RequireInt("runner_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err = h.gl.Runners.DisableProjectRunner(projectID, int64(runnerID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(`{"status": "disabled"}`), nil
}

func (h *handlers) deleteRunner(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runnerID, err := req.RequireInt("runner_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err = h.gl.Runners.RemoveRunner(int64(runnerID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(`{"status": "deleted"}`), nil
}
