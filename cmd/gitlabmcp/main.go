package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyalaletin/gitlabmcp/internal/client"
	"github.com/ilyalaletin/gitlabmcp/internal/config"
	"github.com/ilyalaletin/gitlabmcp/internal/deploy"
	"github.com/ilyalaletin/gitlabmcp/internal/groups"
	"github.com/ilyalaletin/gitlabmcp/internal/issues"
	"github.com/ilyalaletin/gitlabmcp/internal/mr"
	"github.com/ilyalaletin/gitlabmcp/internal/pipelines"
	"github.com/ilyalaletin/gitlabmcp/internal/projects"
	"github.com/ilyalaletin/gitlabmcp/internal/registry"
	"github.com/ilyalaletin/gitlabmcp/internal/releases"
	"github.com/ilyalaletin/gitlabmcp/internal/repositories"
	"github.com/ilyalaletin/gitlabmcp/internal/runners"
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

	// Domain registrations
	deploy.Register(s, gl)
	groups.Register(s, gl)
	issues.Register(s, gl)
	mr.Register(s, gl)
	pipelines.Register(s, gl)
	projects.Register(s, gl)
	registry.Register(s, gl)
	releases.Register(s, gl)
	repositories.Register(s, gl)
	runners.Register(s, gl)

	slog.Info("starting gitlabmcp", "version", version, "gitlab_url", cfg.URL)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
