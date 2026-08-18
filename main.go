package main

import (
	"fmt"
	"math"
)

// QueryResult simulates a Prometheus API query response
type QueryResult struct {
	SeriesCount int
	Limit       int
	IsTruncated bool
}

// BUG: Uses >= instead of > — false positive when result == limit
func checkTruncatedBug(count, limit int) bool {
	return count >= limit // Off-by-one: should be count > limit
}

// FIX: Correct comparison — only truncate when count EXCEEDS limit
func checkTruncatedFixed(count, limit int) bool {
	return count > limit
}

func main() {
	testCases := []struct {
		count int
		limit int
	}{
		{5, 10},   // Less than limit — NOT truncated
		{10, 10},  // Equal to limit — NOT truncated (but bug says it is)
		{15, 10},  // More than limit — IS truncated
		{0, 10},   // Empty result — NOT truncated
		{100, 100},// Equal to limit — NOT truncated (bug triggers here)
	}

	fmt.Println("=== BUGGY VERSION (>= instead of >) ===")
	for _, tc := range testCases {
		bug := checkTruncatedBug(tc.count, tc.limit)
		fmt.Printf("  count=%d limit=%d => truncated=%v %s\n",
			tc.count, tc.limit, bug,
			map[bool]string{true: "⚠️ FALSE POSITIVE", false: ""}[bug && tc.count <= tc.limit])
	}

	fmt.Println("\n=== FIXED VERSION (>) ===")
	for _, tc := range testCases {
		fixed := checkTruncatedFixed(tc.count, tc.limit)
		fmt.Printf("  count=%d limit=%d => truncated=%v\n",
			tc.count, tc.limit, fixed)
	}

	fmt.Println("\nFix: Change `count >= limit` to `count > limit` in the truncation check.")
	_ = math.MaxInt64 // unused import to make go vet happy
}
