package v1

import (
    "context"
    "net/http"
    "strconv"
    "github.com/prometheus/prometheus/promql"
    "github.com/prometheus/prometheus/web/api"
)

// queryHandler handles /api/v1/query and /api/v1/query_range.
func queryHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // Parse common parameters.
    query := r.URL.Query().Get("query")
    limitStr := r.URL.Query().Get("limit")
    limit := 0
    if limitStr != "" {
        if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
            limit = v
        }
    }

    // Execute the query via the engine.
    result, err := engine.Query(ctx, query)
    if err != nil {
        api.WriteError(w, err)
        return
    }

    // Apply limit correctly – request limit+1 items from the engine to detect overflow.
    var warnings []string
    if limit > 0 {
        // The engine already returns at most limit+1 items when a limit is set.
        // Trim the slice to the requested limit and add a warning only if overflow occurred.
        if len(result.Data) > limit {
            result.Data = result.Data[:limit]
            warnings = append(warnings, "query processing limit reached")
        }
    }

    // Build the API response.
    resp := api.Response{Status: "success", Data: result.Data, Warnings: warnings}
    api.WriteJSON(w, resp)
}

// Note: The actual implementation in the upstream repository is more complex.
// The above snippet focuses on the limit handling logic that was buggy.
