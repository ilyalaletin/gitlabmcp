package projects

import (
	"context"
	"fmt"
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

func (h *handlers) listProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	search := req.GetString("search", "")
	owned := req.GetBool("owned", false)
	membership := req.GetBool("membership", false)
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if search != "" {
		opts.Search = gitlab.Ptr(search)
	}
	if owned {
		opts.Owned = gitlab.Ptr(true)
	}
	if membership {
		opts.Membership = gitlab.Ptr(true)
	}

	projects, resp, err := h.gl.Projects.ListProjects(opts)
	if err != nil {
		slog.Error("list_projects failed", "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_projects", "count", len(projects), "duration", time.Since(start))
	result := handler.NewPaginatedResult(projects, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	project, _, err := h.gl.Projects.GetProject(projectID, nil)
	if err != nil {
		slog.Error("get_project failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("get_project", "project_id", projectID, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(project)), nil
}

func (h *handlers) createProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	description := req.GetString("description", "")
	visibility := req.GetString("visibility", "")
	namespaceIDStr := req.GetString("namespace_id", "")

	opts := &gitlab.CreateProjectOptions{
		Name: gitlab.Ptr(name),
	}

	if description != "" {
		opts.Description = gitlab.Ptr(description)
	}
	if visibility != "" {
		opts.Visibility = gitlab.Ptr(gitlab.VisibilityValue(visibility))
	}
	if namespaceIDStr != "" {
		// Parse namespace_id as int64
		var nsID int64
		_, err := fmt.Sscanf(namespaceIDStr, "%d", &nsID)
		if err == nil {
			opts.NamespaceID = gitlab.Ptr(nsID)
		}
	}

	project, _, err := h.gl.Projects.CreateProject(opts)
	if err != nil {
		slog.Error("create_project failed", "name", name, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("create_project", "name", name, "project_id", project.ID, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(project)), nil
}

func (h *handlers) listProjectMembers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectMembersOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	members, resp, err := h.gl.ProjectMembers.ListProjectMembers(projectID, opts)
	if err != nil {
		slog.Error("list_project_members failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_project_members", "project_id", projectID, "count", len(members), "duration", time.Since(start))
	result := handler.NewPaginatedResult(members, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) listProjectWebhooks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectHooksOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	hooks, resp, err := h.gl.Projects.ListProjectHooks(projectID, opts)
	if err != nil {
		slog.Error("list_project_webhooks failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_project_webhooks", "project_id", projectID, "count", len(hooks), "duration", time.Since(start))
	result := handler.NewPaginatedResult(hooks, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) createProjectWebhook(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	pushEvents := req.GetBool("push_events", false)
	mergeRequestsEvents := req.GetBool("merge_requests_events", false)
	issuesEvents := req.GetBool("issues_events", false)
	pipelineEvents := req.GetBool("pipeline_events", false)
	token := req.GetString("token", "")

	opts := &gitlab.AddProjectHookOptions{
		URL: gitlab.Ptr(url),
	}

	if pushEvents {
		opts.PushEvents = gitlab.Ptr(true)
	}
	if mergeRequestsEvents {
		opts.MergeRequestsEvents = gitlab.Ptr(true)
	}
	if issuesEvents {
		opts.IssuesEvents = gitlab.Ptr(true)
	}
	if pipelineEvents {
		opts.PipelineEvents = gitlab.Ptr(true)
	}
	if token != "" {
		opts.Token = gitlab.Ptr(token)
	}

	hook, _, err := h.gl.Projects.AddProjectHook(projectID, opts)
	if err != nil {
		slog.Error("create_project_webhook failed", "project_id", projectID, "url", url, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("create_project_webhook", "project_id", projectID, "webhook_id", hook.ID, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(hook)), nil
}
