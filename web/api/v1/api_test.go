package v1

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "strconv"
)

func TestQueryTruncation_UnderLimit(t *testing.T) {
    // Setup a mock engine that returns 9 series.
    mockEngine := &mockEngine{resultSize: 9}
    // Replace the global engine used by the handler (simplified).
    engine = mockEngine

    req := httptest.NewRequest("GET", "/api/v1/query?query=up&limit=10", nil)
    w := httptest.NewRecorder()
    queryHandler(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
    }
    // Decode JSON and verify warnings are empty.
    var apiResp struct{ Warnings []string }
    json.NewDecoder(resp.Body).Decode(&apiResp)
    if len(apiResp.Warnings) != 0 {
        t.Fatalf("expected no warnings, got %v", apiResp.Warnings)
    }
}

func TestQueryTruncation_ExactLimit(t *testing.T) {
    mockEngine := &mockEngine{resultSize: 10}
    engine = mockEngine

    req := httptest.NewRequest("GET", "/api/v1/query?query=up&limit=10", nil)
    w := httptest.NewRecorder()
    queryHandler(w, req)

    var apiResp struct{ Warnings []string }
    json.NewDecoder(w.Body).Decode(&apiResp)
    if len(apiResp.Warnings) != 0 {
        t.Fatalf("expected no warnings for exact limit, got %v", apiResp.Warnings)
    }
}

func TestQueryTruncation_OverLimit(t *testing.T) {
    mockEngine := &mockEngine{resultSize: 11}
    engine = mockEngine

    req := httptest.NewRequest("GET", "/api/v1/query?query=up&limit=10", nil)
    w := httptest.NewRecorder()
    queryHandler(w, req)

    var apiResp struct{ Warnings []string }
    json.NewDecoder(w.Body).Decode(&apiResp)
    if len(apiResp.Warnings) != 1 {
        t.Fatalf("expected one truncation warning, got %v", apiResp.Warnings)
    }
    if apiResp.Warnings[0] != "query processing limit reached" {
        t.Fatalf("unexpected warning message: %s", apiResp.Warnings[0])
    }
}

// mockEngine implements the minimal subset of the Engine interface needed for tests.
type mockEngine struct {
    resultSize int
}

func (m *mockEngine) Query(ctx context.Context, expr string, limit int) (*Result, error) {
    // Simulate returning `resultSize` series.
    data := make([]interface{}, m.resultSize)
    return &Result{Data: data}, nil
}
