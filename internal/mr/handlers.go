package mr

import (
	"context"
	"log/slog"
	"strconv"
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

func (h *handlers) listMergeRequests(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListProjectMergeRequestsOptions{
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
	if s := req.GetString("author_username", ""); s != "" {
		opts.AuthorUsername = gitlab.Ptr(s)
	}
	if s := req.GetString("reviewer_username", ""); s != "" {
		opts.ReviewerUsername = gitlab.Ptr(s)
	}
	if s := req.GetString("source_branch", ""); s != "" {
		opts.SourceBranch = gitlab.Ptr(s)
	}
	if s := req.GetString("target_branch", ""); s != "" {
		opts.TargetBranch = gitlab.Ptr(s)
	}
	if s := req.GetString("search", ""); s != "" {
		opts.Search = gitlab.Ptr(s)
	}

	mrs, resp, err := h.gl.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		slog.Error("list_merge_requests failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_merge_requests", "project", projectID, "count", len(mrs), "duration", time.Since(start))
	result := handler.NewPaginatedResult(mrs, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getMergeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	mr, _, err := h.gl.MergeRequests.GetMergeRequest(projectID, int64(iid), nil)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(mr)), nil
}

func (h *handlers) createMergeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	sourceBranch, err := req.RequireString("source_branch")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	targetBranch, err := req.RequireString("target_branch")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(title),
		SourceBranch: gitlab.Ptr(sourceBranch),
		TargetBranch: gitlab.Ptr(targetBranch),
	}

	if s := req.GetString("description", ""); s != "" {
		opts.Description = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}
	if s := req.GetString("assignee_ids", ""); s != "" {
		ids := parseInt64IDs(s)
		opts.AssigneeIDs = &ids
	}
	if s := req.GetString("reviewer_ids", ""); s != "" {
		ids := parseInt64IDs(s)
		opts.ReviewerIDs = &ids
	}

	mr, _, err := h.gl.MergeRequests.CreateMergeRequest(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(mr)), nil
}

func (h *handlers) updateMergeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.UpdateMergeRequestOptions{}
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
	if s := req.GetString("target_branch", ""); s != "" {
		opts.TargetBranch = gitlab.Ptr(s)
	}

	mr, _, err := h.gl.MergeRequests.UpdateMergeRequest(projectID, int64(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(mr)), nil
}

func (h *handlers) acceptMergeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.AcceptMergeRequestOptions{}
	if s := req.GetString("merge_commit_message", ""); s != "" {
		opts.MergeCommitMessage = gitlab.Ptr(s)
	}
	opts.Squash = gitlab.Ptr(req.GetBool("squash", false))
	opts.ShouldRemoveSourceBranch = gitlab.Ptr(req.GetBool("should_remove_source_branch", false))

	mr, _, err := h.gl.MergeRequests.AcceptMergeRequest(projectID, int64(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(mr)), nil
}

func (h *handlers) approveMergeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	approval, _, err := h.gl.MergeRequestApprovals.ApproveMergeRequest(projectID, int64(iid), &gitlab.ApproveMergeRequestOptions{})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(approval)), nil
}

func (h *handlers) listMRNotes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	perPage := req.GetInt("per_page", 20)
	perPageClamped := int64(handler.ClampPerPage(int(perPage)))

	opts := &gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(page),
			PerPage: perPageClamped,
		},
	}

	notes, resp, err := h.gl.Notes.ListMergeRequestNotes(projectID, int64(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	result := handler.NewPaginatedResult(notes, int(resp.TotalItems), int(page), int(perPageClamped))
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) createMRNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, _, err := h.gl.Notes.CreateMergeRequestNote(projectID, int64(iid), &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(body),
	})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(note)), nil
}

func (h *handlers) getMRDiff(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("merge_request_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	diffs, _, err := h.gl.MergeRequests.ListMergeRequestDiffs(projectID, int64(iid), &gitlab.ListMergeRequestDiffsOptions{})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(diffs)), nil
}

// parseInt64IDs converts comma-separated string of IDs to []int64 slice
func parseInt64IDs(s string) []int64 {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if id, err := strconv.ParseInt(part, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
