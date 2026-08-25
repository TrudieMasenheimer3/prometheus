package v1

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/teststorage"
)

func TestEvaluateQuery_BoundaryConditions(t *testing.T) {
	ctx := context.Background()

	// Initialize a real TSDB test storage
	storage := teststorage.New(t)

	// Populate the storage with 15 series
	now := time.Now()
	app := storage.Appender(ctx)
	for i := 0; i < 15; i++ {
		lbls := labels.FromStrings("__name__", "http_requests_total", "instance", fmt.Sprintf("inst_%d", i))
		_, err := app.Append(0, lbls, now.UnixMilli(), 1.0)
		if err != nil {
			t.Fatalf("failed to append sample: %v", err)
		}
	}
	if err := app.Commit(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Initialize a real PromQL engine
	reg := prometheus.NewRegistry()
	opts := promql.EngineOpts{
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Reg:        reg,
		MaxSamples: 10000,
		Timeout:    10 * time.Second,
	}
	engine := promql.NewEngine(opts)
	defer engine.Close()

	// Real PromQL engine fetcher for instant query
	instantFetcher := func(limit int) []string {
		qOpts := promql.NewPrometheusQueryOpts(false, 0)
		qry, err := engine.NewInstantQuery(ctx, storage, qOpts, "http_requests_total", now)
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}
		defer qry.Close()

		res := qry.Exec(ctx)
		if res.Err != nil {
			t.Fatalf("failed to execute query: %v", res.Err)
		}

		var series []string
		vector, ok := res.Value.(promql.Vector)
		if !ok {
			t.Fatalf("unexpected type: %T", res.Value)
		}

		for _, sample := range vector {
			series = append(series, sample.Metric.String())
		}

		if limit > 0 && len(series) > limit {
			return series[:limit]
		}
		return series
	}

	// Real PromQL engine fetcher for range query
	rangeFetcher := func(limit int) []string {
		qOpts := promql.NewPrometheusQueryOpts(false, 0)
		start := now.Add(-5 * time.Minute)
		end := now
		step := 1 * time.Minute

		qry, err := engine.NewRangeQuery(ctx, storage, qOpts, "http_requests_total", start, end, step)
		if err != nil {
			t.Fatalf("failed to create range query: %v", err)
		}
		defer qry.Close()

		res := qry.Exec(ctx)
		if res.Err != nil {
			t.Fatalf("failed to execute range query: %v", res.Err)
		}

		var series []string
		matrix, ok := res.Value.(promql.Matrix)
		if !ok {
			t.Fatalf("unexpected type: %T", res.Value)
		}

		for _, s := range matrix {
			series = append(series, s.Metric.String())
		}

		if limit > 0 && len(series) > limit {
			return series[:limit]
		}
		return series
	}

	// Under Limit Test Case
	t.Run("Under Limit - Instant Query", func(t *testing.T) {
		res, err := EvaluateQuery(ctx, instantFetcher, 20)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 15 {
			t.Errorf("expected 15 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 0 {
			t.Errorf("expected 0 warnings, got %d: %v", len(res.Warnings), res.Warnings)
		}
	})

	// Exact Limit Test Case
	t.Run("Exact Limit - Instant Query", func(t *testing.T) {
		res, err := EvaluateQuery(ctx, instantFetcher, 15)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 15 {
			t.Errorf("expected 15 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 0 {
			t.Errorf("expected 0 warnings, got %d: %v", len(res.Warnings), res.Warnings)
		}
	})

	// Over Limit Test Case
	t.Run("Over Limit - Instant Query", func(t *testing.T) {
		res, err := EvaluateQuery(ctx, instantFetcher, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 10 {
			t.Errorf("expected 10 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 1 {
			t.Errorf("expected 1 warning, got %d", len(res.Warnings))
		} else if res.Warnings[0] != TruncationWarning {
			t.Errorf("expected warning %q, got %q", TruncationWarning, res.Warnings[0])
		}
	})

	// Range Query - Under Limit Test Case
	t.Run("Under Limit - Range Query", func(t *testing.T) {
		res, err := EvaluateRangeQuery(ctx, rangeFetcher, 20)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 15 {
			t.Errorf("expected 15 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 0 {
			t.Errorf("expected 0 warnings, got %d: %v", len(res.Warnings), res.Warnings)
		}
	})

	// Range Query - Exact Limit Test Case
	t.Run("Exact Limit - Range Query", func(t *testing.T) {
		res, err := EvaluateRangeQuery(ctx, rangeFetcher, 15)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 15 {
			t.Errorf("expected 15 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 0 {
			t.Errorf("expected 0 warnings, got %d: %v", len(res.Warnings), res.Warnings)
		}
	})

	// Range Query - Over Limit Test Case
	t.Run("Over Limit - Range Query", func(t *testing.T) {
		res, err := EvaluateRangeQuery(ctx, rangeFetcher, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Series) != 10 {
			t.Errorf("expected 10 series, got %d", len(res.Series))
		}
		if len(res.Warnings) != 1 {
			t.Errorf("expected 1 warning, got %d", len(res.Warnings))
		} else if res.Warnings[0] != TruncationWarning {
			t.Errorf("expected warning %q, got %q", TruncationWarning, res.Warnings[0])
		}
	})
}
