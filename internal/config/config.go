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
