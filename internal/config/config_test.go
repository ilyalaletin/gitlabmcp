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
