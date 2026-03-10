# GitLab MCP Server Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an MCP server in Go that exposes 57 tools for interacting with GitLab API over stdio transport.

**Architecture:** Domain-based package structure under `internal/`. Each domain exports `Register(server, client)`. Shared helpers handle parameter extraction, pagination, and error formatting. Config from env vars, stdio transport.

**Tech Stack:** Go, `github.com/mark3labs/mcp-go`, `gitlab.com/gitlab-org/api/client-go`

**Spec:** `docs/superpowers/specs/2026-03-10-gitlab-mcp-server-design.md`

---

## Chunk 1: Foundation

### Task 1: Project Initialization

**Files:**
- Create: `go.mod`
- Create: `Makefile`
- Create: `.gitignore`

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/ilya/dev/gitlabmcp
go mod init github.com/ilya/gitlabmcp
```

- [ ] **Step 2: Add dependencies**

```bash
go get github.com/mark3labs/mcp-go@latest
go get gitlab.com/gitlab-org/api/client-go@latest
```

- [ ] **Step 3: Create .gitignore**

```gitignore
gitlabmcp
*.exe
.env
```

- [ ] **Step 4: Create Makefile**

```makefile
BINARY=gitlabmcp
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test lint clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/gitlabmcp

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
```

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum Makefile .gitignore
git commit -m "feat: initialize Go module with dependencies"
```

---

### Task 2: Config Package

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("GITLAB_TOKEN")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when GITLAB_TOKEN is not set")
	}
}

func TestLoad_WithToken(t *testing.T) {
	t.Setenv("GITLAB_TOKEN", "test-token")
	t.Setenv("GITLAB_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", cfg.Token)
	}
	if cfg.URL != "https://gitlab.com" {
		t.Errorf("expected default URL, got %q", cfg.URL)
	}
}

func TestLoad_CustomURL(t *testing.T) {
	t.Setenv("GITLAB_TOKEN", "test-token")
	t.Setenv("GITLAB_URL", "https://gitlab.example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://gitlab.example.com" {
		t.Errorf("expected custom URL, got %q", cfg.URL)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/config/
```

Expected: FAIL — `Load` not defined.

- [ ] **Step 3: Implement config**

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	Token string
	URL   string
}

func Load() (*Config, error) {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITLAB_TOKEN environment variable is required")
	}

	url := os.Getenv("GITLAB_URL")
	if url == "" {
		url = "https://gitlab.com"
	}

	return &Config{
		Token: token,
		URL:   url,
	}, nil
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/config/ -v
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat: add config package for env var loading"
```

---

### Task 3: Client Wrapper

**Files:**
- Create: `internal/client/client.go`
- Test: `internal/client/client_test.go`

- [ ] **Step 1: Write failing test**

```go
// internal/client/client_test.go
package client

import (
	"testing"

	"github.com/ilya/gitlabmcp/internal/config"
)

func TestNew_DefaultURL(t *testing.T) {
	cfg := &config.Config{
		Token: "test-token",
		URL:   "https://gitlab.com",
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNew_CustomURL(t *testing.T) {
	cfg := &config.Config{
		Token: "test-token",
		URL:   "https://gitlab.example.com",
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/client/
```

Expected: FAIL — `New` not defined.

- [ ] **Step 3: Implement client wrapper**

```go
// internal/client/client.go
package client

import (
	"github.com/ilya/gitlabmcp/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func New(cfg *config.Config) (*gitlab.Client, error) {
	opts := []gitlab.ClientOptionFunc{}
	if cfg.URL != "https://gitlab.com" {
		opts = append(opts, gitlab.WithBaseURL(cfg.URL))
	}
	return gitlab.NewClient(cfg.Token, opts...)
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/client/ -v
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/client/
git commit -m "feat: add GitLab client wrapper"
```

---

### Task 4: Shared Handler Helpers

**Files:**
- Create: `internal/handler/helpers.go`
- Test: `internal/handler/helpers_test.go`

- [ ] **Step 1: Write failing tests**

```go
// internal/handler/helpers_test.go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestPaginationResult(t *testing.T) {
	items := []string{"a", "b", "c"}
	result := NewPaginatedResult(items, 100, 1, 20)

	var out PaginatedResponse[string]
	if err := json.Unmarshal([]byte(result), &out); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if out.TotalCount != 100 {
		t.Errorf("expected total_count 100, got %d", out.TotalCount)
	}
	if out.Page != 1 {
		t.Errorf("expected page 1, got %d", out.Page)
	}
	if out.PerPage != 20 {
		t.Errorf("expected per_page 20, got %d", out.PerPage)
	}
	if len(out.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(out.Items))
	}
}

func TestFormatGitLabError(t *testing.T) {
	msg := FormatGitLabError(errors.New("some error"))
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestClampPerPage(t *testing.T) {
	if ClampPerPage(0) != 20 {
		t.Error("expected default 20 for 0")
	}
	if ClampPerPage(50) != 50 {
		t.Error("expected 50")
	}
	if ClampPerPage(200) != 100 {
		t.Error("expected clamped to 100")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/handler/
```

Expected: FAIL.

- [ ] **Step 3: Implement helpers**

**Important type note:** The GitLab client-go library uses `int64` for all IDs, page numbers, and counts. The mcp-go `CallToolRequest` has built-in typed getters (`req.GetString`, `req.GetInt` which returns `int64`, `req.RequireString`, `req.RequireInt`). Handlers should use these built-in getters directly for parameter extraction. The helpers package only provides pagination formatting, error formatting, and `per_page` clamping.

```go
// internal/handler/helpers.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

func NewPaginatedResult[T any](items []T, totalCount, page, perPage int) string {
	resp := PaginatedResponse[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

// FormatGitLabError extracts status code and message from GitLab API errors.
// For 429 (rate limit), includes Retry-After header value if available.
func FormatGitLabError(err error) string {
	if errResp, ok := err.(*gitlab.ErrorResponse); ok {
		code := errResp.Response.StatusCode
		msg := errResp.Message
		if code == http.StatusTooManyRequests {
			retryAfter := errResp.Response.Header.Get("Retry-After")
			if retryAfter != "" {
				return fmt.Sprintf("GitLab API rate limited (429): %s. Retry after %s seconds.", msg, retryAfter)
			}
		}
		return fmt.Sprintf("GitLab API error (%d): %s", code, msg)
	}
	return fmt.Sprintf("GitLab API error: %s", err.Error())
}

// ClampPerPage enforces default (20) and max (100) for per_page parameter.
func ClampPerPage(perPage int) int {
	if perPage <= 0 {
		return 20
	}
	if perPage > 100 {
		return 100
	}
	return perPage
}

func MarshalJSON(v any) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/handler/ -v
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/handler/
git commit -m "feat: add shared handler helpers for pagination and params"
```

---

### Task 5: Main Entry Point (Skeleton)

**Files:**
- Create: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Create main.go with server skeleton**

```go
// cmd/gitlabmcp/main.go
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilya/gitlabmcp/internal/client"
	"github.com/ilya/gitlabmcp/internal/config"
	"github.com/mark3labs/mcp-go/server"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	gl, err := client.New(cfg)
	if err != nil {
		slog.Error("failed to create GitLab client", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer(
		"gitlabmcp",
		version,
		server.WithToolCapabilities(false),
	)

	// Domain registrations will be added here as each domain is implemented
	_ = gl

	slog.Info("starting gitlabmcp", "version", version, "gitlab_url", cfg.URL)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 2: Verify build**

```bash
make build
```

Expected: binary `gitlabmcp` created without errors.

- [ ] **Step 3: Commit**

```bash
git add cmd/gitlabmcp/main.go
git commit -m "feat: add main entry point with server skeleton"
```

---

## Chunk 2: Core Domains (Issues, Merge Requests, Pipelines)

### Task 6: Issues Domain

**Files:**
- Create: `internal/issues/tools.go`
- Create: `internal/issues/handlers.go`
- Test: `internal/issues/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go` — add `issues.Register(s, gl)`

This task establishes the pattern all subsequent domains will follow. Pay extra attention to the structure.

- [ ] **Step 1: Write tools.go — tool definitions and Register function**

```go
// internal/issues/tools.go
package issues

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func Register(s *server.MCPServer, gl *gitlab.Client) {
	h := &handlers{gl: gl}

	s.AddTool(mcp.NewTool("list_issues",
		mcp.WithDescription("List project issues with filters"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path (e.g. 'owner/repo')")),
		mcp.WithString("state", mcp.Description("Filter by state: opened, closed, all")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("assignee_username", mcp.Description("Filter by assignee username")),
		mcp.WithString("milestone", mcp.Description("Filter by milestone title")),
		mcp.WithString("search", mcp.Description("Search in title and description")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listIssues)

	s.AddTool(mcp.NewTool("get_issue",
		mcp.WithDescription("Get a single issue by IID"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
	), h.getIssue)

	s.AddTool(mcp.NewTool("create_issue",
		mcp.WithDescription("Create a new issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("description", mcp.Description("Issue description")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
		mcp.WithString("assignee_ids", mcp.Description("Comma-separated assignee user IDs")),
		mcp.WithString("milestone_id", mcp.Description("Milestone ID")),
	), h.createIssue)

	s.AddTool(mcp.NewTool("update_issue",
		mcp.WithDescription("Update an existing issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithString("state_event", mcp.Description("State change: close or reopen")),
		mcp.WithString("labels", mcp.Description("Comma-separated label names")),
	), h.updateIssue)

	s.AddTool(mcp.NewTool("delete_issue",
		mcp.WithDescription("Delete an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
	), h.deleteIssue)

	s.AddTool(mcp.NewTool("list_issue_notes",
		mcp.WithDescription("List comments on an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
		mcp.WithNumber("per_page", mcp.Description("Items per page (default 20, max 100)")),
	), h.listIssueNotes)

	s.AddTool(mcp.NewTool("create_issue_note",
		mcp.WithDescription("Add a comment to an issue"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID or path")),
		mcp.WithNumber("issue_iid", mcp.Required(), mcp.Description("Issue IID")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment body (markdown)")),
	), h.createIssueNote)
}
```

- [ ] **Step 2: Write handlers.go**

```go
// internal/issues/handlers.go
package issues

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/ilya/gitlabmcp/internal/handler"
	"github.com/mark3labs/mcp-go/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type handlers struct {
	gl *gitlab.Client
}

// All handlers use mcp-go's built-in typed getters on CallToolRequest:
//   req.RequireString("key") -> (string, error)
//   req.GetString("key", "default") -> string
//   req.RequireInt("key") -> (int64, error)
//   req.GetInt("key", defaultVal) -> int64
//   req.GetBool("key", defaultVal) -> bool
// GitLab client-go uses int64 for all IDs and pagination fields.

func (h *handlers) listIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()

	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := int(req.GetInt("page", 1))
	perPage := handler.ClampPerPage(int(req.GetInt("per_page", 20)))

	opts := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	if s := req.GetString("state", ""); s != "" {
		opts.State = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}
	if s := req.GetString("assignee_username", ""); s != "" {
		opts.AssigneeUsername = gitlab.Ptr(s)
	}
	if s := req.GetString("milestone", ""); s != "" {
		opts.Milestone = gitlab.Ptr(s)
	}
	if s := req.GetString("search", ""); s != "" {
		opts.Search = gitlab.Ptr(s)
	}

	issues, resp, err := h.gl.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		slog.Error("list_issues failed", "project", projectID, "error", err, "duration", time.Since(start))
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	slog.Info("list_issues", "project", projectID, "count", len(issues), "duration", time.Since(start))
	result := handler.NewPaginatedResult(issues, resp.TotalItems, page, perPage)
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) getIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	issue, _, err := h.gl.Issues.GetIssue(projectID, int(iid))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) createIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.CreateIssueOptions{
		Title: gitlab.Ptr(title),
	}
	if s := req.GetString("description", ""); s != "" {
		opts.Description = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}

	issue, _, err := h.gl.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) updateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := &gitlab.UpdateIssueOptions{}
	if s := req.GetString("title", ""); s != "" {
		opts.Title = gitlab.Ptr(s)
	}
	if s := req.GetString("description", ""); s != "" {
		opts.Description = gitlab.Ptr(s)
	}
	if s := req.GetString("state_event", ""); s != "" {
		opts.StateEvent = gitlab.Ptr(s)
	}
	if s := req.GetString("labels", ""); s != "" {
		labels := gitlab.LabelOptions(strings.Split(s, ","))
		opts.Labels = &labels
	}

	issue, _, err := h.gl.Issues.UpdateIssue(projectID, int(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(issue)), nil
}

func (h *handlers) deleteIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err = h.gl.Issues.DeleteIssue(projectID, int(iid))
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(`{"status": "deleted"}`), nil
}

func (h *handlers) listIssueNotes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := int(req.GetInt("page", 1))
	perPage := handler.ClampPerPage(int(req.GetInt("per_page", 20)))

	opts := &gitlab.ListIssueNotesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	notes, resp, err := h.gl.Notes.ListIssueNotes(projectID, int(iid), opts)
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	result := handler.NewPaginatedResult(notes, resp.TotalItems, page, perPage)
	return mcp.NewToolResultText(result), nil
}

func (h *handlers) createIssueNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := req.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iid, err := req.RequireInt("issue_iid")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, _, err := h.gl.Notes.CreateIssueNote(projectID, int(iid), &gitlab.CreateIssueNoteOptions{
		Body: gitlab.Ptr(body),
	})
	if err != nil {
		return mcp.NewToolResultError(handler.FormatGitLabError(err)), nil
	}

	return mcp.NewToolResultText(handler.MarshalJSON(note)), nil
}
```

- [ ] **Step 3: Write handler tests**

```go
// internal/issues/handlers_test.go
package issues

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// Tests verify parameter validation. Full API tests require integration setup.
// Helper to build a CallToolRequest with given arguments.
func makeReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestGetIssue_MissingProjectID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getIssue(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing project_id")
	}
}

func TestGetIssue_MissingIID(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.getIssue(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing issue_iid")
	}
}

func TestCreateIssue_MissingTitle(t *testing.T) {
	h := &handlers{gl: nil}
	result, err := h.createIssue(context.Background(), makeReq(map[string]any{"project_id": "owner/repo"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing title")
	}
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/issues/ -v
```

Expected: PASS (parameter validation tests pass without hitting GitLab API).

- [ ] **Step 5: Register in main.go**

Add import and registration call in `cmd/gitlabmcp/main.go`:

```go
// Add to imports:
"github.com/ilya/gitlabmcp/internal/issues"

// Replace `_ = gl` with:
issues.Register(s, gl)
```

- [ ] **Step 6: Verify build**

```bash
make build
```

Expected: builds successfully.

- [ ] **Step 7: Commit**

```bash
git add internal/issues/ cmd/gitlabmcp/main.go
git commit -m "feat: add issues domain (7 tools)"
```

---

### Task 7: Merge Requests Domain

**Files:**
- Create: `internal/mr/tools.go`
- Create: `internal/mr/handlers.go`
- Test: `internal/mr/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go` — add `mr.Register(s, gl)`

Follow the same pattern as Task 6. Key differences:

- [ ] **Step 1: Write tools.go — 9 tools**

Tools to register:
- `list_merge_requests` — params: `project_id`, `state`, `labels`, `author_username`, `reviewer_username`, `source_branch`, `target_branch`, `search`, `page`, `per_page`
- `get_merge_request` — params: `project_id`, `merge_request_iid`
- `create_merge_request` — params: `project_id`, `title`, `description`, `source_branch` (required), `target_branch` (required), `labels`, `assignee_ids`, `reviewer_ids`
- `update_merge_request` — params: `project_id`, `merge_request_iid`, `title`, `description`, `state_event`, `labels`, `target_branch`
- `accept_merge_request` — params: `project_id`, `merge_request_iid`, `merge_commit_message`, `squash`, `should_remove_source_branch`
- `approve_merge_request` — params: `project_id`, `merge_request_iid`
- `list_mr_notes` — params: `project_id`, `merge_request_iid`, `page`, `per_page`
- `create_mr_note` — params: `project_id`, `merge_request_iid`, `body`
- `get_mr_diff` — params: `project_id`, `merge_request_iid`

GitLab client services:
- `h.gl.MergeRequests.ListProjectMergeRequests(projectID, opts)` → `[]*BasicMergeRequest, *Response, error` (note: returns BasicMergeRequest, not full MergeRequest)
- `h.gl.MergeRequests.GetMergeRequest(projectID, mrIID, nil)` → `*MergeRequest, *Response, error` (third param is `*GetMergeRequestsOptions`, pass `nil`)
- `h.gl.MergeRequests.CreateMergeRequest(projectID, opts)` → `*MergeRequest, *Response, error`
- `h.gl.MergeRequests.UpdateMergeRequest(projectID, mrIID, opts)` → `*MergeRequest, *Response, error`
- `h.gl.MergeRequests.AcceptMergeRequest(projectID, mrIID, opts)` → `*MergeRequest, *Response, error`
- `h.gl.MergeRequestApprovals.ApproveMergeRequest(projectID, mrIID, opts)` → `*MergeRequestApprovals, *Response, error`
- `h.gl.Notes.ListMergeRequestNotes(projectID, mrIID, opts)` → `[]*Note, *Response, error`
- `h.gl.Notes.CreateMergeRequestNote(projectID, mrIID, opts)` → `*Note, *Response, error`
- `h.gl.MergeRequests.ListMergeRequestDiffs(projectID, mrIID, opts)` → `[]*MergeRequestDiff, *Response, error`

All IID/ID parameters are `int` (converted from `int64` via `int(req.RequireInt(...))`).

- [ ] **Step 2: Write handlers.go following issues pattern**

Each handler: extract args using `req.RequireString`/`req.RequireInt`/`req.GetString`/`req.GetInt` → build opts → call GitLab API → format result with `handler.MarshalJSON` or `handler.NewPaginatedResult`. Use `handler.FormatGitLabError(err)` for all API errors. Use `handler.ClampPerPage()` for list tools.

For `accept_merge_request`, use `AcceptMergeRequestOptions` with `MergeCommitMessage`, `Squash`, `ShouldRemoveSourceBranch`.

For `get_mr_diff`, use `h.gl.MergeRequests.ListMergeRequestDiffs(projectID, int(mrIID), opts)` → returns `[]*MergeRequestDiff` with actual diff content.

- [ ] **Step 3: Write parameter validation tests**

Test missing `project_id`, missing `merge_request_iid`, missing required fields on `create_merge_request` (`source_branch`, `target_branch`, `title`).

- [ ] **Step 4: Run tests**

```bash
go test ./internal/mr/ -v
```

- [ ] **Step 5: Register in main.go**

```go
import "github.com/ilya/gitlabmcp/internal/mr"
// in main():
mr.Register(s, gl)
```

- [ ] **Step 6: Verify build and commit**

```bash
make build && make test
git add internal/mr/ cmd/gitlabmcp/main.go
git commit -m "feat: add merge requests domain (9 tools)"
```

---

### Task 8: Pipelines & CI/CD Domain

**Files:**
- Create: `internal/pipelines/tools.go`
- Create: `internal/pipelines/handlers.go`
- Test: `internal/pipelines/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go` — add `pipelines.Register(s, gl)`

- [ ] **Step 1: Write tools.go — 10 tools**

Tools to register:
- `list_pipelines` — params: `project_id`, `status` (running/pending/success/failed/canceled), `ref`, `page`, `per_page`
- `get_pipeline` — params: `project_id`, `pipeline_id`
- `list_pipeline_jobs` — params: `project_id`, `pipeline_id`, `page`, `per_page`
- `get_job_log` — params: `project_id`, `job_id`
- `retry_pipeline` — params: `project_id`, `pipeline_id`
- `cancel_pipeline` — params: `project_id`, `pipeline_id`
- `list_ci_variables` — params: `project_id`, `page`, `per_page`
- `create_ci_variable` — params: `project_id`, `key` (required), `value` (required), `protected`, `masked`, `environment_scope`
- `update_ci_variable` — params: `project_id`, `key` (required), `value`, `protected`, `masked`
- `delete_ci_variable` — params: `project_id`, `key` (required)

GitLab client services:
- `h.gl.Pipelines.ListProjectPipelines(projectID, opts)`
- `h.gl.Pipelines.GetPipeline(projectID, pipelineID)`
- `h.gl.Jobs.ListPipelineJobs(projectID, pipelineID, opts)`
- `h.gl.Jobs.GetTraceFile(projectID, jobID)` — returns `*bytes.Reader, *Response, error`
- `h.gl.Pipelines.RetryPipelineBuild(projectID, pipelineID)`
- `h.gl.Pipelines.CancelPipelineBuild(projectID, pipelineID)`
- `h.gl.ProjectVariables.ListVariables(projectID, opts)`
- `h.gl.ProjectVariables.CreateVariable(projectID, opts)`
- `h.gl.ProjectVariables.UpdateVariable(projectID, key, opts)`
- `h.gl.ProjectVariables.RemoveVariable(projectID, key, nil)` — third param is `*RemoveProjectVariableOptions`, pass `nil`

All IDs are `int` (converted from `int64`). Use `req.RequireString`/`req.RequireInt`/`req.GetString`/`req.GetInt` for parameter extraction. Use `handler.FormatGitLabError(err)` and `handler.ClampPerPage()`.

**Note for `get_job_log`:** The trace file returns `*bytes.Reader`. Read it into a string with `io.ReadAll()`, truncate if >50KB to avoid overwhelming LLM context.

- [ ] **Step 2: Write handlers.go**

- [ ] **Step 3: Write parameter validation tests**

- [ ] **Step 4: Run tests, register in main.go, build, commit**

```bash
go test ./internal/pipelines/ -v
make build && make test
git add internal/pipelines/ cmd/gitlabmcp/main.go
git commit -m "feat: add pipelines & CI/CD domain (10 tools)"
```

---

## Chunk 3: Infrastructure Domains (Runners, Projects, Groups)

### Task 9: Runners Domain

**Files:**
- Create: `internal/runners/tools.go`
- Create: `internal/runners/handlers.go`
- Test: `internal/runners/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 5 tools**

Tools:
- `list_runners` — params: `project_id` (optional — if empty, list all accessible runners), `scope` (active/paused/online/offline), `page`, `per_page`
- `get_runner` — params: `runner_id`
- `enable_runner` — params: `project_id`, `runner_id`
- `disable_runner` — params: `project_id`, `runner_id`
- `delete_runner` — params: `runner_id`

GitLab client:
- `h.gl.Runners.ListRunners(opts)` or `h.gl.Runners.ListProjectRunners(projectID, opts)`
- `h.gl.Runners.GetRunnerDetails(runnerID)`
- `h.gl.Runners.EnableProjectRunner(projectID, opts)` — opts has `RunnerID`
- `h.gl.Runners.DisableProjectRunner(projectID, runnerID)`
- `h.gl.Runners.RemoveRunner(runnerID)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add runners domain (5 tools)"
```

---

### Task 10: Projects Domain

**Files:**
- Create: `internal/projects/tools.go`
- Create: `internal/projects/handlers.go`
- Test: `internal/projects/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 6 tools**

Tools:
- `list_projects` — params: `search`, `owned` (bool), `membership` (bool), `page`, `per_page`
- `get_project` — params: `project_id`
- `create_project` — params: `name` (required), `description`, `visibility` (private/internal/public), `namespace_id`
- `list_project_members` — params: `project_id`, `page`, `per_page`
- `list_project_webhooks` — params: `project_id`, `page`, `per_page`
- `create_project_webhook` — params: `project_id`, `url` (required), `push_events` (bool), `merge_requests_events` (bool), `issues_events` (bool), `pipeline_events` (bool), `token`

GitLab client:
- `h.gl.Projects.ListProjects(opts)`
- `h.gl.Projects.GetProject(projectID, nil)`
- `h.gl.Projects.CreateProject(opts)`
- `h.gl.ProjectMembers.ListProjectMembers(projectID, opts)`
- `h.gl.Projects.ListProjectHooks(projectID, opts)`
- `h.gl.Projects.AddProjectHook(projectID, opts)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add projects domain (6 tools)"
```

---

### Task 11: Groups Domain

**Files:**
- Create: `internal/groups/tools.go`
- Create: `internal/groups/handlers.go`
- Test: `internal/groups/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 4 tools**

Tools:
- `list_groups` — params: `search`, `owned` (bool), `page`, `per_page`
- `get_group` — params: `group_id` (required)
- `list_group_projects` — params: `group_id` (required), `search`, `page`, `per_page`
- `list_group_members` — params: `group_id` (required), `page`, `per_page`

GitLab client:
- `h.gl.Groups.ListGroups(opts)`
- `h.gl.Groups.GetGroup(groupID, nil)`
- `h.gl.Groups.ListGroupProjects(groupID, opts)`
- `h.gl.GroupMembers.ListGroupMembers(groupID, opts)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add groups domain (4 tools)"
```

---

## Chunk 4: Remaining Domains (Repositories, Deploy, Releases, Registry)

### Task 12: Repositories Domain

**Files:**
- Create: `internal/repositories/tools.go`
- Create: `internal/repositories/handlers.go`
- Test: `internal/repositories/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 5 tools**

Tools:
- `list_branches` — params: `project_id`, `search`, `page`, `per_page`
- `get_file` — params: `project_id`, `file_path` (required), `ref` (required, branch/tag/commit SHA)
- `list_repository_tree` — params: `project_id`, `path` (default root), `ref`, `recursive` (bool), `page`, `per_page`
- `list_commits` — params: `project_id`, `ref_name`, `path`, `since` (ISO date), `until` (ISO date), `page`, `per_page`
- `get_commit` — params: `project_id`, `sha` (required)

GitLab client:
- `h.gl.Branches.ListBranches(projectID, opts)`
- `h.gl.RepositoryFiles.GetFile(projectID, filePath, opts)` — opts has `Ref`; response has `Content` (base64)
- `h.gl.Repositories.ListTree(projectID, opts)` — opts has `Path`, `Ref`, `Recursive`
- `h.gl.Commits.ListCommits(projectID, opts)`
- `h.gl.Commits.GetCommit(projectID, sha)`

**Note for `get_file`:** The file content comes base64-encoded. Decode it before returning. Truncate files >100KB.

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add repositories domain (5 tools)"
```

---

### Task 13: Deploy & Environments Domain

**Files:**
- Create: `internal/deploy/tools.go`
- Create: `internal/deploy/handlers.go`
- Test: `internal/deploy/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 5 tools**

Tools:
- `list_environments` — params: `project_id`, `page`, `per_page`
- `get_environment` — params: `project_id`, `environment_id`
- `create_environment` — params: `project_id`, `name` (required), `external_url`
- `stop_environment` — params: `project_id`, `environment_id`
- `list_deployments` — params: `project_id`, `environment`, `status`, `page`, `per_page`

GitLab client:
- `h.gl.Environments.ListEnvironments(projectID, opts)`
- `h.gl.Environments.GetEnvironment(projectID, envID)`
- `h.gl.Environments.CreateEnvironment(projectID, opts)`
- `h.gl.Environments.StopEnvironment(projectID, envID)`
- `h.gl.Deployments.ListProjectDeployments(projectID, opts)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add deploy & environments domain (5 tools)"
```

---

### Task 14: Releases Domain

**Files:**
- Create: `internal/releases/tools.go`
- Create: `internal/releases/handlers.go`
- Test: `internal/releases/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 3 tools**

Tools:
- `list_releases` — params: `project_id`, `page`, `per_page`
- `get_release` — params: `project_id`, `tag_name` (required)
- `create_release` — params: `project_id`, `tag_name` (required), `name`, `description`, `ref` (branch/commit to tag from)

GitLab client:
- `h.gl.Releases.ListReleases(projectID, opts)`
- `h.gl.Releases.GetRelease(projectID, tagName)`
- `h.gl.Releases.CreateRelease(projectID, opts)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add releases domain (3 tools)"
```

---

### Task 15: Container Registry Domain

**Files:**
- Create: `internal/registry/tools.go`
- Create: `internal/registry/handlers.go`
- Test: `internal/registry/handlers_test.go`
- Modify: `cmd/gitlabmcp/main.go`

- [ ] **Step 1: Write tools.go — 3 tools**

Tools:
- `list_registry_repos` — params: `project_id`, `page`, `per_page`
- `list_registry_tags` — params: `project_id`, `repository_id` (required), `page`, `per_page`
- `delete_registry_tag` — params: `project_id`, `repository_id` (required), `tag_name` (required)

GitLab client:
- `h.gl.ContainerRegistry.ListProjectRegistryRepositories(projectID, opts)`
- `h.gl.ContainerRegistry.ListRegistryRepositoryTags(projectID, repoID, opts)`
- `h.gl.ContainerRegistry.DeleteRegistryRepositoryTag(projectID, repoID, tagName)`

- [ ] **Step 2-4: Implement handlers, tests, register, build, commit**

```bash
git commit -m "feat: add container registry domain (3 tools)"
```

---

## Chunk 5: Final Integration

### Task 16: Final Build Verification and Cleanup

**Files:**
- Verify: `cmd/gitlabmcp/main.go` has all domain registrations
- Verify: All tests pass

- [ ] **Step 1: Verify main.go has all 10 domain registrations**

```go
// cmd/gitlabmcp/main.go should have:
issues.Register(s, gl)
mr.Register(s, gl)
pipelines.Register(s, gl)
runners.Register(s, gl)
projects.Register(s, gl)
groups.Register(s, gl)
repositories.Register(s, gl)
deploy.Register(s, gl)
releases.Register(s, gl)
registry.Register(s, gl)
```

- [ ] **Step 2: Run all tests**

```bash
make test
```

Expected: all PASS.

- [ ] **Step 3: Build and verify binary**

```bash
make build
./gitlabmcp --help 2>&1 || true
```

- [ ] **Step 4: Test with MCP inspector (manual)**

```bash
GITLAB_TOKEN=your-token npx @modelcontextprotocol/inspector ./gitlabmcp
```

Verify: all 57 tools appear in the inspector tool list.

- [ ] **Step 5: Final commit**

```bash
git add -A
git commit -m "feat: complete GitLab MCP server with 57 tools"
```
