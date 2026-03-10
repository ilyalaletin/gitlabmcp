package deploy

import (
	"context"
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

func (h *handlers) listEnvironments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListEnvironmentsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	environments, resp, err := h.gl.Environments.ListEnvironments(projectID, opts)
	if err != nil {
		slog.Error("list_environments failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_environments", "project", projectID, "count", len(environments), "duration", time.Since(start))
	result := handler.NewPaginatedResult(environments, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getEnvironment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	environmentID, err := req.RequireInt("environment_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	environment, _, err := h.gl.Environments.GetEnvironment(projectID, int64(environmentID))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(environment)), nil
}

func (h *handlers) createEnvironment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.CreateEnvironmentOptions{
		Name: gitlab.Ptr(name),
	}

	if s := req.GetString("external_url", ""); s != "" {
		opts.ExternalURL = gitlab.Ptr(s)
	}

	environment, _, err := h.gl.Environments.CreateEnvironment(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(environment)), nil
}

func (h *handlers) stopEnvironment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	environmentID, err := req.RequireInt("environment_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, _, err = h.gl.Environments.StopEnvironment(projectID, int64(environmentID), &gitlab.StopEnvironmentOptions{})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText("Environment stopped successfully"), nil
}

func (h *handlers) listDeployments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectDeploymentsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if s := req.GetString("environment", ""); s != "" {
		opts.Environment = gitlab.Ptr(s)
	}
	if s := req.GetString("status", ""); s != "" {
		opts.Status = gitlab.Ptr(s)
	}

	deployments, resp, err := h.gl.Deployments.ListProjectDeployments(projectID, opts)
	if err != nil {
		slog.Error("list_deployments failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_deployments", "project", projectID, "count", len(deployments), "duration", time.Since(start))
	result := handler.NewPaginatedResult(deployments, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}
