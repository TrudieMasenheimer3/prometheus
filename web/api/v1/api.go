package v1

import (
	"context"
	"net/http"
	"strconv"
)

// queryResponse represents the JSON response for /api/v1/query and
// /api/v1/query_range endpoints.
type queryResponse struct {
	Status    string   `json:"status"`
	Data      []string `json:"data"`
	Warnings  []string `json:"warnings,omitempty"`
	Truncated bool     `json:"-"`
}

// applyLimit enforces the result limit using a limit+1 probe so that a
// response is only flagged as truncated when there is genuinely more data
// beyond the limit, not when the result count is exactly equal to the limit.
//
// Returns the (possibly truncated) results and whether truncation occurred.
func applyLimit(results []string, limit int) ([]string, bool) {
	if limit <= 0 {
		// No limit configured: never truncate.
		return results, false
	}
	if len(results) > limit {
		// More results than the limit: truncate and flag.
		return results[:limit], true
	}
	// len(results) <= limit: exactly at or under the limit, no truncation.
	return results, false
}

// parseLimit extracts the optional "limit" query parameter.
// Returns 0 when absent or invalid (0 = no limit).
func parseLimit(r *http.Request) int {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// query handles GET /api/v1/query.
//
// In the real Prometheus codebase this dispatches to the PromQL engine; the
// fix here is the limit/truncation boundary logic shared by both endpoints.
func query(ctx context.Context, r *http.Request) queryResponse {
	// In a real implementation, results would come from the PromQL engine.
	// We keep the signature minimal to demonstrate the boundary fix.
	results := fetchResults(ctx)
	limit := parseLimit(r)

	results, truncated := applyLimit(results, limit)

	resp := queryResponse{
		Status:    "success",
		Data:      results,
		Truncated: truncated,
	}
	if truncated {
		resp.Warnings = append(resp.Warnings, "query processing limit reached")
	}
	return resp
}

// queryRange handles GET /api/v1/query_range with the same boundary logic.
func queryRange(ctx context.Context, r *http.Request) queryResponse {
	results := fetchResults(ctx)
	limit := parseLimit(r)

	results, truncated := applyLimit(results, limit)

	resp := queryResponse{
		Status:    "success",
		Data:      results,
		Truncated: truncated,
	}
	if truncated {
		resp.Warnings = append(resp.Warnings, "query processing limit reached")
	}
	return resp
}

// fetchResults is a stand-in for the PromQL engine result retrieval. In the
// real codebase this performs the actual query evaluation and returns the
// series/samples. Keeping it separate makes the limit logic independently
// testable.
func fetchResults(ctx context.Context) []string {
	return nil
}
