package releases

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_releases",
		mcp.WithDescription("List releases for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listReleases)

	s.AddTool(mcp.NewTool("get_release",
		mcp.WithDescription("Get a single release by tag name"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("tag_name", mcp.Required(), mcp.Description("Release tag name")),
	), h.getRelease)

	s.AddTool(mcp.NewTool("create_release",
		mcp.WithDescription("Create a new release"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("tag_name", mcp.Required(), mcp.Description("Release tag name")),
		mcp.WithString("name", mcp.Description("Release name")),
		mcp.WithString("description", mcp.Description("Release description")),
		mcp.WithString("ref", mcp.Description("Branch, commit SHA, or tag to create the release from")),
	), h.createRelease)
}
