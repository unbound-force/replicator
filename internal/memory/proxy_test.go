package memory

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer creates an httptest server that responds to JSON-RPC calls.
// The handler function receives the method and params and returns a result.
func newTestServer(t *testing.T, handler func(method string, params json.RawMessage) (any, *jsonRPCError)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		paramsBytes, _ := json.Marshal(req.Params)
		result, rpcErr := handler(req.Method, paramsBytes)

		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
		}
		if rpcErr != nil {
			resp.Error = rpcErr
		} else {
			resp.Result, _ = json.Marshal(result)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestCall_Success(t *testing.T) {
	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		return map[string]string{"status": "ok"}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	result, err := client.Call("test_method", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Call: %v", err)
	}

	var parsed map[string]string
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if parsed["status"] != "ok" {
		t.Errorf("status = %q, want %q", parsed["status"], "ok")
	}
}

func TestCall_RPCError(t *testing.T) {
	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		return nil, &jsonRPCError{Code: -32600, Message: "invalid request"}
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Call("test_method", nil)
	if err == nil {
		t.Fatal("expected error for RPC error response")
	}
	if got := err.Error(); got != "dewey error -32600: invalid request" {
		t.Errorf("error = %q, want dewey error message", got)
	}
}

func TestCall_ConnectionRefused(t *testing.T) {
	// Use a URL that will refuse connections.
	client := NewClient("http://127.0.0.1:1")
	_, err := client.Call("test_method", nil)
	if err == nil {
		t.Fatal("expected error for connection refused")
	}

	var unavail *UnavailableError
	if !errors.As(err, &unavail) {
		t.Errorf("expected UnavailableError, got %T: %v", err, err)
	}
}

func TestCall_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Call("test_method", nil)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}

	var unavail *UnavailableError
	if !errors.As(err, &unavail) {
		t.Errorf("expected UnavailableError, got %T: %v", err, err)
	}
}

func TestHealth_Success(t *testing.T) {
	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		if method != "dewey_health" {
			t.Errorf("method = %q, want %q", method, "dewey_health")
		}
		return map[string]string{"status": "healthy"}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	if err := client.Health(); err != nil {
		t.Fatalf("Health: %v", err)
	}
}

func TestHealth_Failure(t *testing.T) {
	client := NewClient("http://127.0.0.1:1")
	err := client.Health()
	if err == nil {
		t.Fatal("expected error for unreachable Dewey")
	}
}

func TestStore_Success(t *testing.T) {
	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		if method != "store_learning" {
			t.Errorf("method = %q, want %q", method, "store_learning")
		}

		var p map[string]string
		json.Unmarshal(params, &p)
		if p["information"] != "test learning" {
			t.Errorf("information = %q, want %q", p["information"], "test learning")
		}
		if p["tags"] != "go,testing" {
			t.Errorf("tags = %q, want %q", p["tags"], "go,testing")
		}

		return map[string]any{"id": "mem-123", "stored": true}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	result, err := client.Store("test learning", "go,testing")
	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	if result["_warning"] == nil {
		t.Error("expected deprecation warning in response")
	}
	warning, ok := result["_warning"].(string)
	if !ok || warning == "" {
		t.Error("expected non-empty deprecation warning string")
	}
}

func TestStore_NoTags(t *testing.T) {
	var receivedParams map[string]any

	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		json.Unmarshal(params, &receivedParams)
		return map[string]any{"stored": true}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Store("info only", "")
	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	if _, hasTags := receivedParams["tags"]; hasTags {
		t.Error("tags should not be sent when empty")
	}
}

func TestStore_DeweyUnavailable(t *testing.T) {
	client := NewClient("http://127.0.0.1:1")
	_, err := client.Store("test", "")
	if err == nil {
		t.Fatal("expected error for unreachable Dewey")
	}

	var unavail *UnavailableError
	if !errors.As(err, &unavail) {
		t.Errorf("expected UnavailableError, got %T", err)
	}
}

func TestFind_Success(t *testing.T) {
	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		if method != "semantic_search" {
			t.Errorf("method = %q, want %q", method, "semantic_search")
		}

		var p map[string]any
		json.Unmarshal(params, &p)
		if p["query"] != "test query" {
			t.Errorf("query = %v, want %q", p["query"], "test query")
		}

		return map[string]any{
			"results": []map[string]string{
				{"page": "test-page", "score": "0.95"},
			},
		}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	result, err := client.Find("test query", "", 5)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	if result["_warning"] == nil {
		t.Error("expected deprecation warning in response")
	}
}

func TestFind_WithCollection(t *testing.T) {
	var receivedParams map[string]any

	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		json.Unmarshal(params, &receivedParams)
		return map[string]any{"results": []any{}}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Find("query", "learnings", 10)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	if receivedParams["source_type"] != "learnings" {
		t.Errorf("source_type = %v, want %q", receivedParams["source_type"], "learnings")
	}
}

func TestFind_WithLimit(t *testing.T) {
	var receivedParams map[string]any

	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		json.Unmarshal(params, &receivedParams)
		return map[string]any{"results": []any{}}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Find("query", "", 7)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	// JSON numbers unmarshal as float64.
	if receivedParams["limit"] != float64(7) {
		t.Errorf("limit = %v, want 7", receivedParams["limit"])
	}
}

func TestFind_ZeroLimit(t *testing.T) {
	var receivedParams map[string]any

	srv := newTestServer(t, func(method string, params json.RawMessage) (any, *jsonRPCError) {
		json.Unmarshal(params, &receivedParams)
		return map[string]any{"results": []any{}}, nil
	})
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Find("query", "", 0)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	if _, hasLimit := receivedParams["limit"]; hasLimit {
		t.Error("limit should not be sent when zero")
	}
}

func TestFind_DeweyUnavailable(t *testing.T) {
	client := NewClient("http://127.0.0.1:1")
	_, err := client.Find("test", "", 5)
	if err == nil {
		t.Fatal("expected error for unreachable Dewey")
	}
}

func TestUnavailableResponse(t *testing.T) {
	err := &UnavailableError{Cause: errors.New("connection refused")}
	resp := UnavailableResponse(err)

	var parsed map[string]any
	if jsonErr := json.Unmarshal([]byte(resp), &parsed); jsonErr != nil {
		t.Fatalf("unmarshal: %v", jsonErr)
	}

	if parsed["code"] != "DEWEY_UNAVAILABLE" {
		t.Errorf("code = %v, want %q", parsed["code"], "DEWEY_UNAVAILABLE")
	}
}

func TestNewClient_Timeout(t *testing.T) {
	client := NewClient("http://example.com")
	if client.http.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", client.http.Timeout)
	}
}
