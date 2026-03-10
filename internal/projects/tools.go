package projects

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_projects",
		mcp.WithDescription("List projects with optional filters"),
		mcp.WithString("search", mcp.Description("Search projects by name or path")),
		mcp.WithBoolean("owned", mcp.Description("Filter projects owned by the authenticated user")),
		mcp.WithBoolean("membership", mcp.Description("Filter projects the user has a membership in")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listProjects)

	s.AddTool(mcp.NewTool("get_project",
		mcp.WithDescription("Get a single project by ID or path"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g., owner/project)")),
	), h.getProject)

	s.AddTool(mcp.NewTool("create_project",
		mcp.WithDescription("Create a new project"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithString("description", mcp.Description("Project description")),
		mcp.WithString("visibility", mcp.Description("Visibility level: private, internal, or public")),
		mcp.WithString("namespace_id", mcp.Description("Namespace ID (group ID) for the project")),
	), h.createProject)

	s.AddTool(mcp.NewTool("list_project_members",
		mcp.WithDescription("List members of a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listProjectMembers)

	s.AddTool(mcp.NewTool("list_project_webhooks",
		mcp.WithDescription("List webhooks for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listProjectWebhooks)

	s.AddTool(mcp.NewTool("create_project_webhook",
		mcp.WithDescription("Create a webhook for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("url", mcp.Required(), mcp.Description("Webhook URL")),
		mcp.WithBoolean("push_events", mcp.Description("Trigger webhook on push events")),
		mcp.WithBoolean("merge_requests_events", mcp.Description("Trigger webhook on merge request events")),
		mcp.WithBoolean("issues_events", mcp.Description("Trigger webhook on issue events")),
		mcp.WithBoolean("pipeline_events", mcp.Description("Trigger webhook on pipeline events")),
		mcp.WithString("token", mcp.Description("Secret token for webhook authentication")),
	), h.createProjectWebhook)
}
