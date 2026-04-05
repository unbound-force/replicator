// Package memory provides a Dewey HTTP proxy client for semantic memory operations.
//
// The hivemind_store and hivemind_find tools proxy to Dewey's semantic search
// endpoints via JSON-RPC 2.0 over HTTP. Six secondary tools return deprecation
// messages pointing users to native Dewey tools.
//
// On connection failure, errors include a structured "DEWEY_UNAVAILABLE" code
// so agents can degrade gracefully.
package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a Dewey HTTP proxy that forwards JSON-RPC calls.
type Client struct {
	url  string
	http *http.Client
}

// NewClient creates a Dewey proxy client with a 10-second timeout.
func NewClient(deweyURL string) *Client {
	return &Client{
		url: deweyURL,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// jsonRPCRequest is a JSON-RPC 2.0 request envelope.
type jsonRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

// jsonRPCResponse is a JSON-RPC 2.0 response envelope.
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
	ID      int             `json:"id"`
}

// jsonRPCError is a JSON-RPC 2.0 error object.
type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Call sends a JSON-RPC 2.0 POST to the Dewey endpoint.
// Returns the result field on success, or a structured error on failure.
func (c *Client) Call(method string, params any) (json.RawMessage, error) {
	reqBody := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.http.Post(c.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, &UnavailableError{Cause: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &UnavailableError{
			Cause: fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody)),
		}
	}

	var rpcResp jsonRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("dewey error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// Health pings the Dewey endpoint to verify connectivity.
func (c *Client) Health() error {
	_, err := c.Call("dewey_health", map[string]any{})
	return err
}

// Store proxies to store_learning with a deprecation warning.
func (c *Client) Store(information, tags string) (map[string]any, error) {
	params := map[string]any{
		"information": information,
	}
	if tags != "" {
		params["tags"] = tags
	}

	result, err := c.Call("store_learning", params)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		// If result isn't a map, wrap it.
		parsed = map[string]any{"result": string(result)}
	}

	// Add deprecation warning to response.
	parsed["_warning"] = "hivemind_store is deprecated. Use dewey_store_learning directly."

	return parsed, nil
}

// Find proxies to semantic_search with a deprecation warning.
func (c *Client) Find(query, collection string, limit int) (map[string]any, error) {
	params := map[string]any{
		"query": query,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	if collection != "" {
		// Collection maps to source_type filter in Dewey.
		params["source_type"] = collection
	}

	result, err := c.Call("semantic_search", params)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		parsed = map[string]any{"result": string(result)}
	}

	// Add deprecation warning to response.
	parsed["_warning"] = "hivemind_find is deprecated. Use dewey_semantic_search directly."

	return parsed, nil
}

// UnavailableError indicates Dewey is not reachable.
type UnavailableError struct {
	Cause error
}

func (e *UnavailableError) Error() string {
	return fmt.Sprintf("dewey unavailable: %v", e.Cause)
}

func (e *UnavailableError) Unwrap() error {
	return e.Cause
}

// UnavailableResponse returns a structured JSON error for agents to parse.
func UnavailableResponse(err error) string {
	resp := map[string]any{
		"error":   err.Error(),
		"code":    "DEWEY_UNAVAILABLE",
		"message": "Dewey semantic search is not available. Memory operations require a running Dewey instance.",
	}
	out, _ := json.MarshalIndent(resp, "", "  ")
	return string(out)
}
