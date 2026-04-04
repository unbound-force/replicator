package mcp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/tools/hive"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

func testServer(t *testing.T) (*Server, *db.Store) {
	t.Helper()
	store, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	reg := registry.New()
	hive.Register(reg, store)
	return NewServer(reg), store
}

func call(t *testing.T, s *Server, method string, params any) json.RawMessage {
	t.Helper()
	paramsJSON, _ := json.Marshal(params)
	req := `{"jsonrpc":"2.0","id":1,"method":"` + method + `","params":` + string(paramsJSON) + "}\n"

	var buf bytes.Buffer
	err := s.Serve(strings.NewReader(req), &buf)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}

	var resp jsonrpcResponse
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v\nraw: %s", err, buf.String())
	}
	if resp.Error != nil {
		t.Fatalf("JSON-RPC error: %d %s", resp.Error.Code, resp.Error.Message)
	}

	result, _ := json.Marshal(resp.Result)
	return result
}

func TestToolsList(t *testing.T) {
	s, _ := testServer(t)
	result := call(t, s, "tools/list", nil)

	var list toolsListResult
	if err := json.Unmarshal(result, &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(list.Tools) != 4 {
		t.Errorf("expected 4 tools, got %d", len(list.Tools))
	}

	names := make(map[string]bool)
	for _, tool := range list.Tools {
		names[tool.Name] = true
	}
	for _, expected := range []string{"hive_cells", "hive_create", "hive_close", "hive_update"} {
		if !names[expected] {
			t.Errorf("missing tool: %s", expected)
		}
	}
}

func TestToolsCall_HiveCells_Empty(t *testing.T) {
	s, _ := testServer(t)
	result := call(t, s, "tools/call", toolsCallParams{
		Name:      "hive_cells",
		Arguments: json.RawMessage(`{}`),
	})

	var callResult toolsCallResult
	json.Unmarshal(result, &callResult)

	if len(callResult.Content) == 0 {
		t.Fatal("expected content in response")
	}
	if callResult.Content[0].Text != "[]" {
		t.Errorf("expected empty array, got %s", callResult.Content[0].Text)
	}
}

func TestToolsCall_HiveCreate(t *testing.T) {
	s, _ := testServer(t)
	result := call(t, s, "tools/call", toolsCallParams{
		Name:      "hive_create",
		Arguments: json.RawMessage(`{"title": "Test cell", "type": "bug"}`),
	})

	var callResult toolsCallResult
	json.Unmarshal(result, &callResult)

	if len(callResult.Content) == 0 {
		t.Fatal("expected content")
	}

	var cell map[string]any
	json.Unmarshal([]byte(callResult.Content[0].Text), &cell)

	if cell["title"] != "Test cell" {
		t.Errorf("title = %v, want %q", cell["title"], "Test cell")
	}
	if cell["type"] != "bug" {
		t.Errorf("type = %v, want %q", cell["type"], "bug")
	}
}

func TestToolsCall_CreateThenQuery(t *testing.T) {
	s, _ := testServer(t)

	// Create a cell.
	call(t, s, "tools/call", toolsCallParams{
		Name:      "hive_create",
		Arguments: json.RawMessage(`{"title": "My task"}`),
	})

	// Query cells.
	result := call(t, s, "tools/call", toolsCallParams{
		Name:      "hive_cells",
		Arguments: json.RawMessage(`{}`),
	})

	var callResult toolsCallResult
	json.Unmarshal(result, &callResult)

	var cells []map[string]any
	json.Unmarshal([]byte(callResult.Content[0].Text), &cells)

	if len(cells) != 1 {
		t.Fatalf("expected 1 cell, got %d", len(cells))
	}
	if cells[0]["title"] != "My task" {
		t.Errorf("title = %v, want %q", cells[0]["title"], "My task")
	}
}

func TestToolsCall_UnknownTool(t *testing.T) {
	s, _ := testServer(t)

	req := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nonexistent","arguments":{}}}` + "\n"
	var buf bytes.Buffer
	s.Serve(strings.NewReader(req), &buf)

	var resp jsonrpcResponse
	json.Unmarshal(buf.Bytes(), &resp)

	if resp.Error == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestInitialize(t *testing.T) {
	s, _ := testServer(t)
	result := call(t, s, "initialize", map[string]any{})

	var initResult map[string]any
	json.Unmarshal(result, &initResult)

	if initResult["protocolVersion"] != "2024-11-05" {
		t.Errorf("protocolVersion = %v", initResult["protocolVersion"])
	}
}
