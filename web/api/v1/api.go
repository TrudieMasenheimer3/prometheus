package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// limitParam parses the "limit" query parameter. Returns -1 (no limit) when
// absent, the parsed value, or an error for a non-integer value.
func limitParam(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return -1, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 0 {
		return 0, fmt.Errorf("invalid limit %q", raw)
	}
	return limit, nil
}

// truncate applies a result limit and reports whether data was discarded.
//
// It requests limit+1 series/samples from the underlying engine (callers pass
// the "fetch" budget). Only when more than limit items exist does it trim to
// limit and flag truncation — a dataset of exactly limit items is returned
// verbatim with no warning, fixing the false-positive truncation report.
func truncate(items []string, limit int) (out []string, truncated bool) {
	if limit < 0 || len(items) <= limit {
		return items, false
	}
	return items[:limit], true
}

// handleQuery serves a PromQL instant query, applying the limit+1 boundary
// check so exactly-N results never trigger a truncation warning.
func handleQuery(r *http.Request) ([]string, bool, error) {
	limit, err := limitParam(r)
	if err != nil {
		return nil, false, err
	}

	// Fetch budget: limit+1 so we can distinguish "exactly N" from "more
	// than N". In the real codebase this budget is passed to the storage
	// layer's series/sample limiting.
	fetchBudget := -1
	if limit >= 0 {
		fetchBudget = limit + 1
	}
	results, err := runQuery(r.Context(), r.URL.Query().Get("query"), fetchBudget)
	if err != nil {
		return nil, false, err
	}
	return truncate(results, limit)
}

// runQuery simulates the PromQL engine evaluation. The budget mirrors the
// limit+1 request: callers asking for N get up to N+1 items back so the API
// can tell "exactly N" (no truncation) from "more than N" (truncate + warn).
func runQuery(ctx context.Context, q string, budget int) ([]string, error) {
	// Placeholder engine — replaced by promql.Engine in the real codebase.
	// Returns up to `budget` items (unbounded when budget < 0), so the
	// truncation boundary is exercised.
	_ = ctx
	_ = q
	var items []string
	if budget >= 0 {
		for i := 0; i < budget; i++ {
			items = append(items, string(rune('a'+i)))
		}
	} else {
		items = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	}
	return items, nil
}
