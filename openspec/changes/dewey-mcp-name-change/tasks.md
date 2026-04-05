## 1. Update Dewey proxy tool names

- [x] 1.1 In `internal/memory/proxy.go`, change `c.Call("dewey_dewey_store_learning", ...)` to `c.Call("store_learning", ...)`
- [x] 1.2 In `internal/memory/proxy.go`, change `c.Call("dewey_dewey_semantic_search", ...)` to `c.Call("semantic_search", ...)`
- [x] 1.3 In `internal/memory/proxy.go`, update the `Store` comment from `dewey_dewey_store_learning` to `store_learning`
- [x] 1.4 In `internal/memory/proxy.go`, update the `Find` comment from `dewey_dewey_semantic_search` to `semantic_search`

## 2. Update deprecation warnings

- [x] 2.1 In `internal/memory/proxy.go`, change the Store `_warning` text from `"Use dewey_dewey_store_learning directly."` to `"Use dewey_store_learning directly."`
- [x] 2.2 In `internal/memory/proxy.go`, change the Find `_warning` text from `"Use dewey_dewey_semantic_search directly."` to `"Use dewey_semantic_search directly."`

## 3. Update tests

- [x] 3.1 In `internal/memory/proxy_test.go`, update the Store test assertion from `"dewey_dewey_store_learning"` to `"store_learning"`
- [x] 3.2 In `internal/memory/proxy_test.go`, update the Find test assertion from `"dewey_dewey_semantic_search"` to `"semantic_search"`

## 4. Verify

- [x] 4.1 Run `go vet ./...` -- zero warnings
- [x] 4.2 Run `go test ./internal/memory/ -count=1 -v` -- all proxy tests pass with new tool names
- [x] 4.3 Run `go test ./... -count=1` -- full suite passes, no regressions
