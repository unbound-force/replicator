package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/memory"
	commstools "github.com/unbound-force/replicator/internal/tools/comms"
	forgetools "github.com/unbound-force/replicator/internal/tools/forge"
	memorytools "github.com/unbound-force/replicator/internal/tools/memory"
	"github.com/unbound-force/replicator/internal/tools/org"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

func buildFullRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	store, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	reg := registry.New()
	org.Register(reg, store)
	commstools.Register(reg, store)
	forgetools.Register(reg, store)
	memClient := memory.NewClient("http://localhost:3333/mcp/")
	memorytools.Register(reg, memClient)
	return reg
}

func TestWriteDocs_ContainsAllTools(t *testing.T) {
	reg := buildFullRegistry(t)

	var buf bytes.Buffer
	if err := writeDocs(&buf, reg); err != nil {
		t.Fatalf("writeDocs: %v", err)
	}

	output := buf.String()

	// Verify all registered tools appear in the output.
	for _, tool := range reg.List() {
		if !strings.Contains(output, tool.Name) {
			t.Errorf("tool %q not found in docs output", tool.Name)
		}
	}
}

func TestWriteDocs_HasCategoryHeaders(t *testing.T) {
	reg := buildFullRegistry(t)

	var buf bytes.Buffer
	writeDocs(&buf, reg)
	output := buf.String()

	for _, header := range []string{"## Org", "## Comms", "## Forge", "## Memory"} {
		if !strings.Contains(output, header) {
			t.Errorf("missing category header: %q", header)
		}
	}
}

func TestWriteDocs_ToolCount(t *testing.T) {
	reg := buildFullRegistry(t)

	if reg.Count() < 50 {
		t.Errorf("expected at least 50 tools, got %d", reg.Count())
	}

	var buf bytes.Buffer
	writeDocs(&buf, reg)
	output := buf.String()

	if !strings.Contains(output, "tools registered") {
		t.Error("output missing tool count line")
	}
}
