package registry

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_registry_repos",
		mcp.WithDescription("List container registry repositories for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listRegistryRepos)

	s.AddTool(mcp.NewTool("list_registry_tags",
		mcp.WithDescription("List tags in a container registry repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("repository_id", mcp.Required(), mcp.Description("Repository ID")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listRegistryTags)

	s.AddTool(mcp.NewTool("delete_registry_tag",
		mcp.WithDescription("Delete a tag from a container registry repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("repository_id", mcp.Required(), mcp.Description("Repository ID")),
		mcp.WithString("tag_name", mcp.Required(), mcp.Description("Tag name to delete")),
	), h.deleteRegistryTag)
}
