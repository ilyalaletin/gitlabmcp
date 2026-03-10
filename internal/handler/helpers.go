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
