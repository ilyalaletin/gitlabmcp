package pipelines

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_pipelines",
		mcp.WithDescription("List project pipelines with optional filters"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g. 'owner/repo')")),
		mcp.WithString("status", mcp.Description("Filter by status: running, pending, success, failed, canceled")),
		mcp.WithString("ref", mcp.Description("Filter by ref (branch/tag)")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listPipelines)

	s.AddTool(mcp.NewTool("get_pipeline",
		mcp.WithDescription("Get a single pipeline by ID"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("pipeline_id", mcp.Required(), mcp.Description("Pipeline ID")),
	), h.getPipeline)

	s.AddTool(mcp.NewTool("list_pipeline_jobs",
		mcp.WithDescription("List jobs in a pipeline"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("pipeline_id", mcp.Required(), mcp.Description("Pipeline ID")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listPipelineJobs)

	s.AddTool(mcp.NewTool("get_job_log",
		mcp.WithDescription("Get job log output (truncated to 50KB if larger)"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("job_id", mcp.Required(), mcp.Description("Job ID")),
	), h.getJobLog)

	s.AddTool(mcp.NewTool("retry_pipeline",
		mcp.WithDescription("Retry a failed pipeline"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("pipeline_id", mcp.Required(), mcp.Description("Pipeline ID")),
	), h.retryPipeline)

	s.AddTool(mcp.NewTool("cancel_pipeline",
		mcp.WithDescription("Cancel a running pipeline"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("pipeline_id", mcp.Required(), mcp.Description("Pipeline ID")),
	), h.cancelPipeline)

	s.AddTool(mcp.NewTool("list_ci_variables",
		mcp.WithDescription("List CI/CD variables for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listCIVariables)

	s.AddTool(mcp.NewTool("create_ci_variable",
		mcp.WithDescription("Create a CI/CD variable"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Variable key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
		mcp.WithBoolean("protected", mcp.Description("Protect the variable (default false)")),
		mcp.WithBoolean("masked", mcp.Description("Mask the variable in logs (default false)")),
		mcp.WithString("environment_scope", mcp.Description("Environment scope (default '*')")),
	), h.createCIVariable)

	s.AddTool(mcp.NewTool("update_ci_variable",
		mcp.WithDescription("Update a CI/CD variable"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Variable key")),
		mcp.WithString("value", mcp.Description("New variable value")),
		mcp.WithBoolean("protected", mcp.Description("Update protected status")),
		mcp.WithBoolean("masked", mcp.Description("Update masked status")),
	), h.updateCIVariable)

	s.AddTool(mcp.NewTool("delete_ci_variable",
		mcp.WithDescription("Delete a CI/CD variable"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Variable key")),
	), h.deleteCIVariable)
}
