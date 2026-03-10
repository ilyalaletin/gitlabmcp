package handler

import (
	"encoding/json"
	"errors"
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
