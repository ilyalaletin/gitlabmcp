package issues

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_issues",
		mcp.WithDescription("List project issues with filters"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g. 'owner/repo')")),
		mcp.WithString("state", mcp.Description("Filter by state: opened, closed, all")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("assignee_username", mcp.Description("Filter by assignee username")),
		mcp.WithString("milestone", mcp.Description("Filter by milestone title")),
		mcp.WithString("search", mcp.Description("Search in title and description")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listIssues)

	s.AddTool(mcp.NewTool("get_issue",
		mcp.WithDescription("Get a single issue by IID"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
	), h.getIssue)

	s.AddTool(mcp.NewTool("create_issue",
		mcp.WithDescription("Create a new issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("description", mcp.Description("Issue description")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("assignee_ids", mcp.Description("Comma-separated assignee user IDs")),
		mcp.WithString("milestone_id", mcp.Description("Milestone ID")),
	), h.createIssue)

	s.AddTool(mcp.NewTool("update_issue",
		mcp.WithDescription("Update an existing issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithString("state_event", mcp.Description("State change: close or reopen")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
	), h.updateIssue)

	s.AddTool(mcp.NewTool("delete_issue",
		mcp.WithDescription("Delete an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
	), h.deleteIssue)

	s.AddTool(mcp.NewTool("list_issue_notes",
		mcp.WithDescription("List comments on an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listIssueNotes)

	s.AddTool(mcp.NewTool("create_issue_note",
		mcp.WithDescription("Add a comment to an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment body (markdown)")),
	), h.createIssueNote)
}
