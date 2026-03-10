package groups

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

func (h *handlers) listGroups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	search := req.GetString("search", "")
	owned := req.GetBool("owned", false)
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListGroupsOptions{
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

	groups, resp, err := h.gl.Groups.ListGroups(opts)
	if err != nil {
		slog.Error("list_groups failed", "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_groups", "count", len(groups), "duration", time.Since(start))
	result := handler.NewPaginatedResult(groups, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getGroup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID, err := req.RequireString("group_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	group, _, err := h.gl.Groups.GetGroup(groupID, nil)
	if err != nil {
		slog.Error("get_group failed", "group_id", groupID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("get_group", "group_id", groupID, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(group)), nil
}

func (h *handlers) listGroupProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID, err := req.RequireString("group_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	search := req.GetString("search", "")
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if search != "" {
		opts.Search = gitlab.Ptr(search)
	}

	projects, resp, err := h.gl.Groups.ListGroupProjects(groupID, opts)
	if err != nil {
		slog.Error("list_group_projects failed", "group_id", groupID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_group_projects", "group_id", groupID, "count", len(projects), "duration", time.Since(start))
	result := handler.NewPaginatedResult(projects, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) listGroupMembers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID, err := req.RequireString("group_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListGroupMembersOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	members, resp, err := h.gl.Groups.ListGroupMembers(groupID, opts)
	if err != nil {
		slog.Error("list_group_members failed", "group_id", groupID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_group_members", "group_id", groupID, "count", len(members), "duration", time.Since(start))
	result := handler.NewPaginatedResult(members, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}
