package deploy

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_environments",
		mcp.WithDescription("List environments in a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g. 'owner/repo')")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listEnvironments)

	s.AddTool(mcp.NewTool("get_environment",
		mcp.WithDescription("Get a single environment by ID"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("environment_id", mcp.Required(), mcp.Description("Environment ID")),
	), h.getEnvironment)

	s.AddTool(mcp.NewTool("create_environment",
		mcp.WithDescription("Create a new environment"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Environment name")),
		mcp.WithString("external_url", mcp.Description("External URL for the environment")),
	), h.createEnvironment)

	s.AddTool(mcp.NewTool("stop_environment",
		mcp.WithDescription("Stop an environment"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("environment_id", mcp.Required(), mcp.Description("Environment ID")),
	), h.stopEnvironment)

	s.AddTool(mcp.NewTool("list_deployments",
		mcp.WithDescription("List deployments in a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("environment", mcp.Description("Filter by environment name")),
		mcp.WithString("status", mcp.Description("Filter by status: created, running, success, failed, canceled")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listDeployments)
}
