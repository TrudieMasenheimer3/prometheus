package v1

import (
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
// The caller requests limit+1 items from the engine; only when more than
// limit items exist does this trim to limit and flag truncation. A dataset of
// exactly limit items is returned verbatim with no warning — fixing the
// false-positive truncation report (the off-by-one at len == limit).
func truncate(items []string, limit int) (out []string, truncated bool) {
	if limit < 0 || len(items) <= limit {
		return items, false
	}
	return items[:limit], true
}
