package mr

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_merge_requests",
		mcp.WithDescription("List project merge requests with filters"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g. 'owner/repo')")),
		mcp.WithString("state", mcp.Description("Filter by state: opened, closed, merged, all")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("author_username", mcp.Description("Filter by author username")),
		mcp.WithString("reviewer_username", mcp.Description("Filter by reviewer username")),
		mcp.WithString("source_branch", mcp.Description("Filter by source branch")),
		mcp.WithString("target_branch", mcp.Description("Filter by target branch")),
		mcp.WithString("search", mcp.Description("Search in title and description")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listMergeRequests)

	s.AddTool(mcp.NewTool("get_merge_request",
		mcp.WithDescription("Get a single merge request by IID"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
	), h.getMergeRequest)

	s.AddTool(mcp.NewTool("create_merge_request",
		mcp.WithDescription("Create a new merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Merge request title")),
		mcp.WithString("description", mcp.Description("Merge request description")),
		mcp.WithString("source_branch", mcp.Required(), mcp.Description("Source branch name")),
		mcp.WithString("target_branch", mcp.Required(), mcp.Description("Target branch name")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("assignee_ids", mcp.Description("Comma-separated assignee user IDs")),
		mcp.WithString("reviewer_ids", mcp.Description("Comma-separated reviewer user IDs")),
	), h.createMergeRequest)

	s.AddTool(mcp.NewTool("update_merge_request",
		mcp.WithDescription("Update an existing merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithString("state_event", mcp.Description("State change: close or reopen")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("target_branch", mcp.Description("New target branch")),
	), h.updateMergeRequest)

	s.AddTool(mcp.NewTool("accept_merge_request",
		mcp.WithDescription("Accept (merge) a merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
		mcp.WithString("merge_commit_message", mcp.Description("Custom merge commit message")),
		mcp.WithBoolean("squash", mcp.Description("Squash commits before merging (default false)")),
		mcp.WithBoolean("should_remove_source_branch", mcp.Description("Delete source branch after merge (default false)")),
	), h.acceptMergeRequest)

	s.AddTool(mcp.NewTool("approve_merge_request",
		mcp.WithDescription("Approve a merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
	), h.approveMergeRequest)

	s.AddTool(mcp.NewTool("list_mr_notes",
		mcp.WithDescription("List comments on a merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listMRNotes)

	s.AddTool(mcp.NewTool("create_mr_note",
		mcp.WithDescription("Add a comment to a merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment body (markdown)")),
	), h.createMRNote)

	s.AddTool(mcp.NewTool("get_mr_diff",
		mcp.WithDescription("Get changes in a merge request"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("merge_request_iid", mcp.Required(), mcp.Description("Merge request IID")),
	), h.getMRDiff)
}
