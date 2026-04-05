package hive

import "testing"

func TestCreateEpic(t *testing.T) {
	store := testStore(t)

	epic, subtasks, err := CreateEpic(store, CreateEpicInput{
		EpicTitle:       "Build the thing",
		EpicDescription: "A big feature",
		Subtasks: []SubtaskInput{
			{Title: "Step 1", Priority: 2},
			{Title: "Step 2", Priority: 1, Files: []string{"foo.go", "bar.go"}},
		},
	})
	if err != nil {
		t.Fatalf("CreateEpic: %v", err)
	}

	if epic.Title != "Build the thing" {
		t.Errorf("epic title = %q, want %q", epic.Title, "Build the thing")
	}
	if epic.Type != "epic" {
		t.Errorf("epic type = %q, want %q", epic.Type, "epic")
	}
	if epic.Status != "open" {
		t.Errorf("epic status = %q, want %q", epic.Status, "open")
	}

	if len(subtasks) != 2 {
		t.Fatalf("expected 2 subtasks, got %d", len(subtasks))
	}
	if subtasks[0].Title != "Step 1" {
		t.Errorf("subtask[0] title = %q, want %q", subtasks[0].Title, "Step 1")
	}
	if subtasks[0].Priority != 2 {
		t.Errorf("subtask[0] priority = %d, want %d", subtasks[0].Priority, 2)
	}
	if subtasks[0].ParentID == nil || *subtasks[0].ParentID != epic.ID {
		t.Errorf("subtask[0] parent_id = %v, want %q", subtasks[0].ParentID, epic.ID)
	}
	if subtasks[1].Title != "Step 2" {
		t.Errorf("subtask[1] title = %q, want %q", subtasks[1].Title, "Step 2")
	}
}

func TestCreateEpic_NoSubtasks(t *testing.T) {
	store := testStore(t)

	epic, subtasks, err := CreateEpic(store, CreateEpicInput{
		EpicTitle: "Empty epic",
	})
	if err != nil {
		t.Fatalf("CreateEpic: %v", err)
	}
	if epic.Title != "Empty epic" {
		t.Errorf("epic title = %q, want %q", epic.Title, "Empty epic")
	}
	if len(subtasks) != 0 {
		t.Errorf("expected 0 subtasks, got %d", len(subtasks))
	}
}

func TestCreateEpic_SubtasksQueryable(t *testing.T) {
	store := testStore(t)

	epic, _, err := CreateEpic(store, CreateEpicInput{
		EpicTitle: "Queryable epic",
		Subtasks: []SubtaskInput{
			{Title: "Sub A"},
			{Title: "Sub B"},
		},
	})
	if err != nil {
		t.Fatalf("CreateEpic: %v", err)
	}

	// Query all cells -- should find epic + 2 subtasks.
	cells, err := QueryCells(store, CellQuery{})
	if err != nil {
		t.Fatalf("QueryCells: %v", err)
	}
	if len(cells) != 3 {
		t.Errorf("expected 3 cells (1 epic + 2 subtasks), got %d", len(cells))
	}

	// Query by epic type.
	epics, err := QueryCells(store, CellQuery{Type: "epic"})
	if err != nil {
		t.Fatalf("QueryCells: %v", err)
	}
	if len(epics) != 1 {
		t.Errorf("expected 1 epic, got %d", len(epics))
	}
	if epics[0].ID != epic.ID {
		t.Errorf("epic ID = %q, want %q", epics[0].ID, epic.ID)
	}
}

func TestCreateEpic_DefaultPriority(t *testing.T) {
	store := testStore(t)

	_, subtasks, err := CreateEpic(store, CreateEpicInput{
		EpicTitle: "Priority test",
		Subtasks: []SubtaskInput{
			{Title: "No priority set"},
		},
	})
	if err != nil {
		t.Fatalf("CreateEpic: %v", err)
	}
	if subtasks[0].Priority != 1 {
		t.Errorf("default priority = %d, want 1", subtasks[0].Priority)
	}
}
