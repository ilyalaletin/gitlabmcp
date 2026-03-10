# GitLab MCP Server — Design Spec

## Overview

MCP server in Go that provides tools for interacting with GitLab API via stdio transport. Covers issues, merge requests, pipelines, runners, projects, deploy, releases, and container registry.

## Tech Stack

- **Language:** Go
- **MCP SDK:** `github.com/mark3labs/mcp-go`
- **GitLab client:** `gitlab.com/gitlab-org/api/client-go`
- **Transport:** stdio

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITLAB_TOKEN` | yes | — | Personal Access Token |
| `GITLAB_URL` | no | `https://gitlab.com` | Base URL of GitLab instance |

Server exits with an error if `GITLAB_TOKEN` is not set.

## Project Structure

```
gitlabmcp/
├── cmd/
│   └── gitlabmcp/
│       └── main.go              # entry point, server init
├── internal/
│   ├── config/
│   │   └── config.go            # env var loading
│   ├── client/
│   │   └── client.go            # gitlab client wrapper
│   ├── issues/
│   │   ├── tools.go             # tool registration
│   │   └── handlers.go          # handlers
│   ├── mr/
│   │   ├── tools.go
│   │   └── handlers.go
│   ├── pipelines/
│   │   ├── tools.go
│   │   └── handlers.go
│   ├── runners/
│   │   ├── tools.go
│   │   └── handlers.go
│   ├── projects/
│   │   ├── tools.go
│   │   └── handlers.go
│   ├── deploy/
│   │   ├── tools.go
│   │   └── handlers.go
│   ├── releases/
│   │   ├── tools.go
│   │   └── handlers.go
│   └── registry/
│       ├── tools.go
│       └── handlers.go
├── go.mod
└── go.sum
```

### Key principles

- `config` reads `GITLAB_TOKEN` and `GITLAB_URL` from environment
- `client` creates `*gitlab.Client` once, shared across all domains
- Each domain package exports `Register(server, client)` to register its tools
- `main.go`: init config -> client -> MCP server -> call Register for each domain -> start stdio

## Tools

### Issues (7 tools)

| Tool | Description |
|------|-------------|
| `list_issues` | List project/group issues with filters (state, labels, assignee, milestone) |
| `get_issue` | Get issue by ID |
| `create_issue` | Create issue |
| `update_issue` | Update issue (title, description, labels, assignee, milestone, state) |
| `delete_issue` | Delete issue |
| `list_issue_notes` | List issue comments |
| `create_issue_note` | Add comment to issue |

### Merge Requests (9 tools)

| Tool | Description |
|------|-------------|
| `list_merge_requests` | List MRs with filters (state, labels, author, reviewer) |
| `get_merge_request` | Get MR by ID |
| `create_merge_request` | Create MR |
| `update_merge_request` | Update MR |
| `merge_merge_request` | Merge MR |
| `approve_merge_request` | Approve MR |
| `list_mr_notes` | List MR comments |
| `create_mr_note` | Add comment to MR |
| `get_mr_diff` | Get MR diff |

### Pipelines & CI/CD (10 tools)

| Tool | Description |
|------|-------------|
| `list_pipelines` | List project pipelines |
| `get_pipeline` | Get pipeline by ID |
| `list_pipeline_jobs` | List jobs in pipeline |
| `get_job_log` | Get job log |
| `retry_pipeline` | Retry pipeline |
| `cancel_pipeline` | Cancel pipeline |
| `list_ci_variables` | List project CI/CD variables |
| `create_ci_variable` | Create CI/CD variable |
| `update_ci_variable` | Update CI/CD variable |
| `delete_ci_variable` | Delete CI/CD variable |

### Runners (5 tools)

| Tool | Description |
|------|-------------|
| `list_runners` | List runners (project/group/global) |
| `get_runner` | Get runner by ID |
| `enable_runner` | Enable runner for project |
| `disable_runner` | Disable runner for project |
| `delete_runner` | Delete runner |

### Projects & Groups (6 tools)

| Tool | Description |
|------|-------------|
| `list_projects` | List projects with search |
| `get_project` | Get project |
| `create_project` | Create project |
| `list_project_members` | List project members |
| `list_project_webhooks` | List webhooks |
| `create_project_webhook` | Create webhook |

### Deploy & Environments (5 tools)

| Tool | Description |
|------|-------------|
| `list_environments` | List project environments |
| `get_environment` | Get environment |
| `create_environment` | Create environment |
| `stop_environment` | Stop environment |
| `list_deployments` | List deployments |

### Releases (3 tools)

| Tool | Description |
|------|-------------|
| `list_releases` | List releases |
| `get_release` | Get release |
| `create_release` | Create release |

### Container Registry (3 tools)

| Tool | Description |
|------|-------------|
| `list_registry_repos` | List registry repositories |
| `list_registry_tags` | List tags |
| `delete_registry_tag` | Delete tag |

**Total: ~45 tools.** All tools take `project_id` (string, `owner/repo` or numeric ID) as required parameter, except global ones like `list_runners`.

## Error Handling

- GitLab API errors (401, 403, 404, 429) are translated to readable MCP error responses
- Rate limiting (429): return error with "rate limited, retry later" message, no automatic retries
- Invalid parameters: validation at handler level with clear messages ("project_id is required")

## Response Format

- Tools return JSON — structured data convenient for LLM consumption
- Lists include pagination: `items` + `total_count` + `page` + `per_page`
- Default `per_page` = 20 to avoid overloading LLM context
