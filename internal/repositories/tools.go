package repositories

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_branches",
		mcp.WithDescription("List branches in a repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("search", mcp.Description("Search branches by name")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listBranches)

	s.AddTool(mcp.NewTool("get_file",
		mcp.WithDescription("Get file content from a repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Path to the file in the repository")),
		mcp.WithString("ref", mcp.Required(), mcp.Description("Branch, tag, or commit SHA")),
	), h.getFile)

	s.AddTool(mcp.NewTool("list_repository_tree",
		mcp.WithDescription("List files and directories in a repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("path", mcp.Description("Path to list (default /)")),
		mcp.WithString("ref", mcp.Description("Branch, tag, or commit SHA (default main/master)")),
		mcp.WithBoolean("recursive", mcp.Description("Recursively list all files")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listRepositoryTree)

	s.AddTool(mcp.NewTool("list_commits",
		mcp.WithDescription("List commits in a repository"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("ref_name", mcp.Description("Branch, tag, or commit SHA")),
		mcp.WithString("path", mcp.Description("Filter commits by file path")),
		mcp.WithString("since", mcp.Description("ISO 8601 date (e.g., 2024-01-01T00:00:00Z)")),
		mcp.WithString("until", mcp.Description("ISO 8601 date (e.g., 2024-01-31T23:59:59Z)")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listCommits)

	s.AddTool(mcp.NewTool("get_commit",
		mcp.WithDescription("Get a single commit"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("sha", mcp.Required(), mcp.Description("Commit SHA")),
	), h.getCommit)
}
