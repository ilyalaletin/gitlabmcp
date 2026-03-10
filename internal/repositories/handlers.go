package repositories

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/ilya/gitlabmcp/internal/handler"
	"github.com/mark3labs/mcp-go/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// parseISO8601 parses an ISO 8601 date string and returns a *time.Time pointer or nil if empty
func parseISO8601(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try parsing without time component
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil
		}
	}
	return &t
}

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

func (h *handlers) listBranches(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	search := req.GetString("search", "")
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if search != "" {
		opts.Search = gitlab.Ptr(search)
	}

	branches, resp, err := h.gl.Branches.ListBranches(projectID, opts)
	if err != nil {
		slog.Error("list_branches failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_branches", "project_id", projectID, "count", len(branches), "duration", time.Since(start))
	result := handler.NewPaginatedResult(branches, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getFile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	filePath, err := req.RequireString("file_path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ref, err := req.RequireString("ref")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	opts := &gitlab.GetFileOptions{
		Ref: gitlab.Ptr(ref),
	}

	file, _, err := h.gl.RepositoryFiles.GetFile(projectID, filePath, opts)
	if err != nil {
		slog.Error("get_file failed", "project_id", projectID, "file_path", filePath, "ref", ref, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	// File content is base64-encoded; decode it
	decoded, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		slog.Error("base64 decode failed", "file_path", filePath, "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to decode file content: %v", err)), nil
	}

	// Truncate if larger than 100KB
	content := string(decoded)
	const maxSize = 100 * 1024
	if len(content) > maxSize {
		content = content[:maxSize] + "\n... (truncated)"
	}

	slog.Info("get_file", "project_id", projectID, "file_path", filePath, "ref", ref, "size", len(decoded), "duration", time.Since(start))
	return mcp.NewToolResultText(content), nil
}

func (h *handlers) listRepositoryTree(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	path := req.GetString("path", "/")
	ref := req.GetString("ref", "")
	recursive := req.GetBool("recursive", false)
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListTreeOptions{
		Path:      gitlab.Ptr(path),
		Recursive: gitlab.Ptr(recursive),
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if ref != "" {
		opts.Ref = gitlab.Ptr(ref)
	}

	nodes, resp, err := h.gl.Repositories.ListTree(projectID, opts)
	if err != nil {
		slog.Error("list_repository_tree failed", "project_id", projectID, "path", path, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_repository_tree", "project_id", projectID, "path", path, "count", len(nodes), "duration", time.Since(start))
	result := handler.NewPaginatedResult(nodes, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) listCommits(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	refName := req.GetString("ref_name", "")
	path := req.GetString("path", "")
	since := req.GetString("since", "")
	until := req.GetString("until", "")
	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if refName != "" {
		opts.RefName = gitlab.Ptr(refName)
	}
	if path != "" {
		opts.Path = gitlab.Ptr(path)
	}
	if sinceTime := parseISO8601(since); sinceTime != nil {
		opts.Since = sinceTime
	}
	if untilTime := parseISO8601(until); untilTime != nil {
		opts.Until = untilTime
	}

	commits, resp, err := h.gl.Commits.ListCommits(projectID, opts)
	if err != nil {
		slog.Error("list_commits failed", "project_id", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_commits", "project_id", projectID, "count", len(commits), "duration", time.Since(start))
	result := handler.NewPaginatedResult(commits, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getCommit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	sha, err := req.RequireString("sha")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start := time.Now()

	commit, _, err := h.gl.Commits.GetCommit(projectID, sha, nil)
	if err != nil {
		slog.Error("get_commit failed", "project_id", projectID, "sha", sha, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("get_commit", "project_id", projectID, "sha", sha, "duration", time.Since(start))
	return mcp.NewToolResultText(handler.MarshalJSON(commit)), nil
}
