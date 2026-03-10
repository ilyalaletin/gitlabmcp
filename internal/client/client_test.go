package client

import (
	"testing"

	"github.com/ilyalaletin/gitlabmcp/internal/config"
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
