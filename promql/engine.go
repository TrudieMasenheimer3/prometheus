package promql

import (
    "context"
    "github.com/prometheus/prometheus/storage"
)

// Query executes a PromQL expression and respects the optional result limit.
func (e *Engine) Query(ctx context.Context, expr string, limit int) (*Result, error) {
    // Parse the expression (omitted for brevity).
    // ...

    // When a limit is set we ask the storage for limit+1 samples/series.
    // This extra sample allows us to know whether the result was truncated.
    fetchLimit := 0
    if limit > 0 {
        fetchLimit = limit + 1
    }

    // Retrieve data from the storage.
    raw, err := e.storage.Select(ctx, storage.SelectParams{Limit: fetchLimit})
    if err != nil {
        return nil, err
    }

    // Convert raw storage data to Result format.
    res := convert(raw)
    return res, nil
}

// convert transforms storage data into the Result type used by the API.
func convert(raw storage.SeriesSet) *Result {
    // Implementation omitted – unchanged from upstream.
    return &Result{}
}

// Note: The real engine contains many more functions; only the limit handling part is shown.
