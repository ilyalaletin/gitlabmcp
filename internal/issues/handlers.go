package issues

import (
	"context"
	"log/slog"
	"strings"
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

func (h *handlers) listIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	if s := req.GetString("state", ""); s != "" {
		opts.State = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}
	if s := req.GetString("assignee_username", ""); s != "" {
		opts.AssigneeUsername = gitlab.Ptr(s)
	}
	if s := req.GetString("milestone", ""); s != "" {
		opts.Milestone = gitlab.Ptr(s)
	}
	if s := req.GetString("search", ""); s != "" {
		opts.Search = gitlab.Ptr(s)
	}

	issues, resp, err := h.gl.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		slog.Error("list_issues failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_issues", "project", projectID, "count", len(issues), "duration", time.Since(start))
	result := handler.NewPaginatedResult(issues, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	issue, _, err := h.gl.Issues.GetIssue(projectID, int64(iid))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) createIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.CreateIssueOptions{
		Title: gitlab.Ptr(title),
	}
	if s := req.GetString("description", ""); s != "" {
		opts.Description = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}

	issue, _, err := h.gl.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) updateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.UpdateIssueOptions{}
	if s := req.GetString("title", ""); s != "" {
		opts.Title = gitlab.Ptr(s)
	}
	if s := req.GetString("description", ""); s != "" {
		opts.Description = gitlab.Ptr(s)
	}
	if s := req.GetString("state_event", ""); s != "" {
		opts.StateEvent = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}

	issue, _, err := h.gl.Issues.UpdateIssue(projectID, int64(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) deleteIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err = h.gl.Issues.DeleteIssue(projectID, int64(iid))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(`{"status": "deleted"}`), nil
}

func (h *handlers) listIssueNotes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListIssueNotesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	notes, resp, err := h.gl.Notes.ListIssueNotes(projectID, int64(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	result := handler.NewPaginatedResult(notes, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) createIssueNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, _, err := h.gl.Notes.CreateIssueNote(projectID, int64(iid), &gitlab.CreateIssueNoteOptions{
		Body: gitlab.Ptr(body),
	})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(note)), nil
}
