package v1

import (
	"context"
	"errors"
)

// TruncationWarning is the standard warning returned when results exceed limit.
const TruncationWarning = "query processing limit reached"

// QueryResult represents a vector or matrix query result.
type QueryResult struct {
	Series   []string
	Warnings []string
}

// EvaluateQuery processes an instant query with an optional limit.
// If limit > 0, it fetches up to limit + 1 items to determine if truncation actually occurred.
func EvaluateQuery(ctx context.Context, fetcher func(limit int) []string, limit int) (*QueryResult, error) {
	if fetcher == nil {
		return nil, errors.New("nil query fetcher")
	}

	if limit <= 0 {
		series := fetcher(0)
		return &QueryResult{Series: series, Warnings: nil}, nil
	}

	// Fetch limit + 1 items to detect if results exceed limit without false positives
	fetchLimit := limit + 1
	fetched := fetcher(fetchLimit)

	var warnings []string
	resultSeries := fetched

	if len(fetched) > limit {
		resultSeries = fetched[:limit]
		warnings = append(warnings, TruncationWarning)
	}

	return &QueryResult{
		Series:   resultSeries,
		Warnings: warnings,
	}, nil
}

// EvaluateRangeQuery processes a range query with an optional limit.
func EvaluateRangeQuery(ctx context.Context, fetcher func(limit int) []string, limit int) (*QueryResult, error) {
	return EvaluateQuery(ctx, fetcher, limit)
}
