package runners

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_runners",
		mcp.WithDescription("List runners with optional filters"),
		mcp.WithString("project_id", mcp.Description("Project ID or path (optional - if empty, list all accessible runners)")),
		mcp.WithString("scope", mcp.Description("Filter by scope: active, paused, online, offline")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listRunners)

	s.AddTool(mcp.NewTool("get_runner",
		mcp.WithDescription("Get a single runner by ID"),
		mcp.WithNumber("runner_id", mcp.Required(), mcp.Description("Runner ID")),
	), h.getRunner)

	s.AddTool(mcp.NewTool("enable_runner",
		mcp.WithDescription("Enable a runner for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("runner_id", mcp.Required(), mcp.Description("Runner ID")),
	), h.enableRunner)

	s.AddTool(mcp.NewTool("disable_runner",
		mcp.WithDescription("Disable a runner for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("runner_id", mcp.Required(), mcp.Description("Runner ID")),
	), h.disableRunner)

	s.AddTool(mcp.NewTool("delete_runner",
		mcp.WithDescription("Delete a runner"),
		mcp.WithNumber("runner_id", mcp.Required(), mcp.Description("Runner ID")),
	), h.deleteRunner)
}
