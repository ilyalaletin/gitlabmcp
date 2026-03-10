package groups

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_groups",
		mcp.WithDescription("List groups with optional filters"),
		mcp.WithString("search", mcp.Description("Search groups by name or path")),
		mcp.WithBoolean("owned", mcp.Description("Filter groups owned by the authenticated user")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listGroups)

	s.AddTool(mcp.NewTool("get_group",
		mcp.WithDescription("Get a single group by ID or path"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID or path (e.g., my-group)")),
	), h.getGroup)

	s.AddTool(mcp.NewTool("list_group_projects",
		mcp.WithDescription("List projects in a group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID or path")),
		mcp.WithString("search", mcp.Description("Search projects by name")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listGroupProjects)

	s.AddTool(mcp.NewTool("list_group_members",
		mcp.WithDescription("List members of a group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listGroupMembers)
}
