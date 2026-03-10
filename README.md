# gitlabmcp

[![CI](https://img.shields.io/github/actions/workflow/status/ilya/gitlabmcp/ci.yml?branch=main&style=flat-square&label=CI)](https://github.com/ilya/gitlabmcp/actions/workflows/ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/ilya/gitlabmcp?style=flat-square)](https://github.com/ilya/gitlabmcp/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.24-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/github/license/ilya/gitlabmcp?style=flat-square)](LICENSE)

An [MCP](https://modelcontextprotocol.io/) server for GitLab, providing 57 tools across 10 domains for seamless GitLab integration with Claude, Claude Desktop, Cursor, and other MCP-compatible clients.

## Quick Start

### Prerequisites

- **GitLab Personal Access Token (PAT)**: Create one at `https://gitlab.com/-/user_settings/personal_access_tokens` with `api` and `read_api` scopes.
- **MCP-compatible client**: Claude Code, Claude Desktop, Cursor, Cline, or similar.

### Installation

**Option 1: Pre-built binary**

Download the latest release from [GitHub Releases](https://github.com/ilya/gitlabmcp/releases):

```bash
# macOS/Linux
wget https://github.com/ilya/gitlabmcp/releases/latest/download/gitlabmcp_v0.1.0_darwin_amd64.tar.gz
tar xzf gitlabmcp_v0.1.0_darwin_amd64.tar.gz
chmod +x gitlabmcp
mv gitlabmcp /usr/local/bin/
```

**Option 2: Build from source**

```bash
go install github.com/ilya/gitlabmcp/cmd/gitlabmcp@latest
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITLAB_TOKEN` | Yes | — | GitLab Personal Access Token |
| `GITLAB_URL` | No | `https://gitlab.com` | GitLab API base URL (for self-hosted instances) |

## MCP Client Setup

### Claude Code

Add the server to your Claude Code configuration:

```bash
claude mcp add gitlab -e GITLAB_TOKEN=glpat-xxx -- gitlabmcp
```

Or manually edit your MCP configuration file.

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "gitlabmcp",
      "env": {
        "GITLAB_TOKEN": "glpat-xxx"
      }
    }
  }
}
```

### Cursor

Create or edit `.cursor/mcp.json` in your workspace root:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "gitlabmcp",
      "env": {
        "GITLAB_TOKEN": "glpat-xxx"
      }
    }
  }
}
```

### Cline / RooCode (VS Code)

Open Extensions → Cline → Settings → MCP Servers, then add:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "gitlabmcp",
      "env": {
        "GITLAB_TOKEN": "glpat-xxx"
      }
    }
  }
}
```

## Tools Reference

### Projects (11 tools)

| Tool | Description |
|------|-------------|
| `list_projects` | List all accessible projects |
| `get_project` | Get project details by ID or path |
| `create_project` | Create a new project |
| `update_project` | Update project settings |
| `delete_project` | Delete a project |
| `archive_project` | Archive a project |
| `unarchive_project` | Unarchive a project |
| `get_project_members` | List project members |
| `add_project_member` | Add a member to a project |
| `update_project_member` | Update member access level |
| `remove_project_member` | Remove a member from a project |

### Groups (4 tools)

| Tool | Description |
|------|-------------|
| `list_groups` | List all groups |
| `get_group` | Get group details |
| `create_group` | Create a new group |
| `update_group` | Update group settings |

### Repositories (5 tools)

| Tool | Description |
|------|-------------|
| `list_branches` | List project branches |
| `get_branch` | Get branch details |
| `get_commit` | Get commit details |
| `list_tags` | List project tags |
| `get_tree` | Get repository tree at path |

### Merge Requests (11 tools)

| Tool | Description |
|------|-------------|
| `list_merge_requests` | List merge requests in project |
| `get_merge_request` | Get MR details |
| `create_merge_request` | Create a new merge request |
| `update_merge_request` | Update MR title, description, etc. |
| `merge_merge_request` | Merge an MR |
| `close_merge_request` | Close an MR |
| `reopen_merge_request` | Reopen a closed MR |
| `approve_merge_request` | Approve an MR |
| `unapprove_merge_request` | Revoke approval |
| `list_merge_request_notes` | List MR comments |
| `create_merge_request_note` | Add a comment to an MR |

### Issues (10 tools)

| Tool | Description |
|------|-------------|
| `list_issues` | List project issues |
| `get_issue` | Get issue details |
| `create_issue` | Create a new issue |
| `update_issue` | Update issue title, description, etc. |
| `close_issue` | Close an issue |
| `reopen_issue` | Reopen a closed issue |
| `list_issue_notes` | List issue comments |
| `create_issue_note` | Add a comment to an issue |
| `add_issue_label` | Add a label to an issue |
| `remove_issue_label` | Remove a label from an issue |

### CI/CD Pipeline (6 tools)

| Tool | Description |
|------|-------------|
| `list_pipelines` | List project pipelines |
| `get_pipeline` | Get pipeline details and status |
| `list_pipeline_jobs` | List jobs in a pipeline |
| `get_job` | Get job details and logs |
| `retry_pipeline` | Retry a pipeline |
| `cancel_pipeline` | Cancel a running pipeline |

### Releases (3 tools)

| Tool | Description |
|------|-------------|
| `list_releases` | List project releases |
| `get_release` | Get release details |
| `create_release` | Create a new release |

### Container Registry (3 tools)

| Tool | Description |
|------|-------------|
| `list_registry_repositories` | List container registry repositories |
| `get_registry_repository` | Get repository details |
| `list_registry_repository_tags` | List tags in a repository |

### Environments & Deployments (5 tools)

| Tool | Description |
|------|-------------|
| `list_environments` | List project environments |
| `get_environment` | Get environment details |
| `list_deployments` | List deployments |
| `get_deployment` | Get deployment details |
| `list_deployment_merge_requests` | List MRs deployed in an environment |

### Users (3 tools)

| Tool | Description |
|------|-------------|
| `get_current_user` | Get authenticated user details |
| `list_users` | List all users (admin only) |
| `get_user` | Get user details by ID |

## Development

### Build

```bash
go build -o gitlabmcp ./cmd/gitlabmcp
```

### Test

```bash
go test ./...
```

### Lint

```bash
go vet ./...
```

### Run Locally

```bash
GITLAB_TOKEN=glpat-xxx ./gitlabmcp
```

## License

This project is licensed under the [MIT License](LICENSE).
