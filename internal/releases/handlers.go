package releases

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

func (h *handlers) listReleases(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListReleasesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	releases, resp, err := h.gl.Releases.ListReleases(projectID, opts)
	if err != nil {
		slog.Error("list_releases failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_releases", "project_id", projectID, "count", len(releases), "duration", time.Since(start))
	result := handler.NewPaginatedResult(releases, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getRelease(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tagName, err := req.RequireString("tag_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	release, _, err := h.gl.Releases.GetRelease(projectID, tagName)
	if err != nil {
		slog.Error("get_release failed", "project_id", projectID, "tag_name", tagName, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("get_release", "project_id", projectID, "tag_name", tagName, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(release)), nil
}

func (h *handlers) createRelease(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tagName, err := req.RequireString("tag_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	opts := &gitlab.CreateReleaseOptions{
		TagName: gitlab.Ptr(tagName),
	}

	if name := req.GetString("name", ""); name != "" {
		opts.Name = gitlab.Ptr(name)
	}

	if description := req.GetString("description", ""); description != "" {
		opts.Description = gitlab.Ptr(description)
	}

	if ref := req.GetString("ref", ""); ref != "" {
		opts.Ref = gitlab.Ptr(ref)
	}

	release, _, err := h.gl.Releases.CreateRelease(projectID, opts)
	if err != nil {
		slog.Error("create_release failed", "project_id", projectID, "tag_name", tagName, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("create_release", "project_id", projectID, "tag_name", tagName, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(release)), nil
}
