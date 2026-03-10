package client

import (
	"github.com/ilyalaletin/gitlabmcp/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func New(cfg *config.Config) (*gitlab.Client, error) {
	opts := []gitlab.ClientOptionFunc{}
	if cfg.URL != "https://gitlab.com" {
		opts = append(opts, gitlab.WithBaseURL(cfg.URL))
	}
	return gitlab.NewClient(cfg.Token, opts...)
}
