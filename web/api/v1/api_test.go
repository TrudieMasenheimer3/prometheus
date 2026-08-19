package v1

import (
	"context"
	"fmt"
	"testing"
)

func generateDummySeries(n int) []string {
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = fmt.Sprintf("http_requests_total{instance=\"inst_%d\"}", i)
	}
	return res
}

func mockFetcher(totalCount int) func(limit int) []string {
	all := generateDummySeries(totalCount)
	return func(limit int) []string {
		if limit <= 0 || limit >= len(all) {
			return all
		}
		return all[:limit]
	}
}

func TestEvaluateQueryUnderLimit(t *testing.T) {
	ctx := context.Background()
	fetcher := mockFetcher(9)
	res, err := EvaluateQuery(ctx, fetcher, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Series) != 9 {
		t.Fatalf("expected 9 series, got %d", len(res.Series))
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("expected 0 warnings, got %v", res.Warnings)
	}
}

func TestEvaluateQueryExactLimit(t *testing.T) {
	ctx := context.Background()
	fetcher := mockFetcher(10)
	res, err := EvaluateQuery(ctx, fetcher, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Series) != 10 {
		t.Fatalf("expected 10 series, got %d", len(res.Series))
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("expected 0 warnings for exact limit boundary, got %v", res.Warnings)
	}
}

func TestEvaluateQueryOverLimit(t *testing.T) {
	ctx := context.Background()
	fetcher := mockFetcher(15)
	res, err := EvaluateQuery(ctx, fetcher, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Series) != 10 {
		t.Fatalf("expected truncated 10 series, got %d", len(res.Series))
	}
	if len(res.Warnings) != 1 || res.Warnings[0] != TruncationWarning {
		t.Fatalf("expected truncation warning, got %v", res.Warnings)
	}
}

func TestEvaluateRangeQueryExactLimit(t *testing.T) {
	ctx := context.Background()
	fetcher := mockFetcher(5)
	res, err := EvaluateRangeQuery(ctx, fetcher, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Series) != 5 {
		t.Fatalf("expected 5 series, got %d", len(res.Series))
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("expected 0 warnings for range query exact limit, got %v", res.Warnings)
	}
}

func TestEvaluateRangeQueryOverLimit(t *testing.T) {
	ctx := context.Background()
	fetcher := mockFetcher(8)
	res, err := EvaluateRangeQuery(ctx, fetcher, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Series) != 5 {
		t.Fatalf("expected 5 series, got %d", len(res.Series))
	}
	if len(res.Warnings) != 1 || res.Warnings[0] != TruncationWarning {
		t.Fatalf("expected truncation warning, got %v", res.Warnings)
	}
}
