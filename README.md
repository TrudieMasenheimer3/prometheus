# Prometheus Query Truncation Fix

Fixes the false positive truncation warning in `/api/v1/query` when the result count equals the limit.

## The Bug

When `count >= limit` is used instead of `count > limit`, a query returning exactly the limit number of results incorrectly reports as truncated.

## The Fix

Change the comparison from `>=` to `>` so truncation is only reported when the result count **exceeds** the limit.

```diff
- if len(results) >= limit {
+ if len(results) > limit {
```

## Files

- `main.go` — Demonstrates both the buggy and fixed behavior with test cases
