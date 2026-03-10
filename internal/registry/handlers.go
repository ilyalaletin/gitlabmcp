package registry

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

func (h *handlers) listRegistryRepos(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectRegistryRepositoriesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	repos, resp, err := h.gl.ContainerRegistry.ListProjectRegistryRepositories(projectID, opts)
	if err != nil {
		slog.Error("list_registry_repos failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_registry_repos", "project_id", projectID, "count", len(repos), "duration", time.Since(start))
	result := handler.NewPaginatedResult(repos, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) listRegistryTags(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	repositoryIDRaw, err := req.RequireInt("repository_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListRegistryRepositoryTagsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	tags, resp, err := h.gl.ContainerRegistry.ListRegistryRepositoryTags(projectID, int64(repositoryIDRaw), opts)
	if err != nil {
		slog.Error("list_registry_tags failed", "project_id", projectID, "repository_id", repositoryIDRaw, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_registry_tags", "project_id", projectID, "repository_id", repositoryIDRaw, "count", len(tags), "duration", time.Since(start))
	result := handler.NewPaginatedResult(tags, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) deleteRegistryTag(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	repositoryIDRaw, err := req.RequireInt("repository_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tagName, err := req.RequireString("tag_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	_, err = h.gl.ContainerRegistry.DeleteRegistryRepositoryTag(projectID, int64(repositoryIDRaw), tagName)
	if err != nil {
		slog.Error("delete_registry_tag failed", "project_id", projectID, "repository_id", repositoryIDRaw, "tag_name", tagName, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("delete_registry_tag", "project_id", projectID, "repository_id", repositoryIDRaw, "tag_name", tagName, "duration", time.Since(start))
	return mcp.NewToolResultText("Registry tag deleted successfully"), nil
}
