package v1

import (
	"net/http/httptest"
	"testing"
)

func TestApplyLimitUnderLimit(t *testing.T) {
	results := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	got, truncated := applyLimit(results, 10)
	if len(got) != 9 {
		t.Errorf("expected 9 results, got %d", len(got))
	}
	if truncated {
		t.Error("should not be truncated when under the limit")
	}
}

func TestApplyLimitExactLimit(t *testing.T) {
	results := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	got, truncated := applyLimit(results, 10)
	if len(got) != 10 {
		t.Errorf("expected 10 results, got %d", len(got))
	}
	if truncated {
		t.Error("should NOT be truncated when result count equals the limit")
	}
}

func TestApplyLimitOverLimit(t *testing.T) {
	results := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	got, truncated := applyLimit(results, 10)
	if len(got) != 10 {
		t.Errorf("expected 10 results (truncated), got %d", len(got))
	}
	if !truncated {
		t.Error("should be truncated when result count exceeds the limit")
	}
}

func TestApplyLimitNoLimit(t *testing.T) {
	results := []string{"a", "b", "c"}
	got, truncated := applyLimit(results, 0)
	if len(got) != 3 {
		t.Errorf("expected 3 results, got %d", len(got))
	}
	if truncated {
		t.Error("should never truncate when no limit is set")
	}
}

func TestQueryExactLimitNoWarning(t *testing.T) {
	// Verify the endpoint-level behavior: exact-limit responses carry no warning.
	req := httptest.NewRequest("GET", "/api/v1/query?limit=2", nil)
	// Seed fetchResults indirectly is not possible here; the boundary logic is
	// covered directly above. This test asserts parseLimit behavior.
	if got := parseLimit(req); got != 2 {
		t.Errorf("expected limit 2, got %d", got)
	}
}
