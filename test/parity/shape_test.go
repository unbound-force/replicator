package parity

import (
	"encoding/json"
	"testing"
)

func TestShapeMatch_IdenticalObjects(t *testing.T) {
	expected := json.RawMessage(`{"name": "alice", "age": 30, "active": true}`)
	actual := json.RawMessage(`{"name": "bob", "age": 25, "active": false}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_TypeMismatch(t *testing.T) {
	expected := json.RawMessage(`{"count": 42}`)
	actual := json.RawMessage(`{"count": "forty-two"}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch, got match")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].Path != "$.count" {
		t.Errorf("expected path $.count, got %s", diffs[0].Path)
	}
	if diffs[0].ExpectedType != "number" {
		t.Errorf("expected type number, got %s", diffs[0].ExpectedType)
	}
	if diffs[0].ActualType != "string" {
		t.Errorf("expected actual type string, got %s", diffs[0].ActualType)
	}
}

func TestShapeMatch_MissingField(t *testing.T) {
	expected := json.RawMessage(`{"id": "abc", "title": "hello", "status": "open"}`)
	actual := json.RawMessage(`{"id": "xyz", "title": "world"}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch, got match")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].Path != "$.status" {
		t.Errorf("expected path $.status, got %s", diffs[0].Path)
	}
	if diffs[0].ActualType != "missing" {
		t.Errorf("expected actual type missing, got %s", diffs[0].ActualType)
	}
}

func TestShapeMatch_ExtraField(t *testing.T) {
	// Extra fields in actual should still match -- compatible superset.
	expected := json.RawMessage(`{"id": "abc"}`)
	actual := json.RawMessage(`{"id": "xyz", "extra_field": "bonus", "another": 42}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match (extra fields OK), got diffs: %+v", diffs)
	}
}

func TestShapeMatch_NestedDifference(t *testing.T) {
	expected := json.RawMessage(`{"user": {"name": "alice", "profile": {"bio": "hello"}}}`)
	actual := json.RawMessage(`{"user": {"name": "bob", "profile": {"bio": 42}}}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch, got match")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].Path != "$.user.profile.bio" {
		t.Errorf("expected path $.user.profile.bio, got %s", diffs[0].Path)
	}
	if diffs[0].ExpectedType != "string" {
		t.Errorf("expected type string, got %s", diffs[0].ExpectedType)
	}
	if diffs[0].ActualType != "number" {
		t.Errorf("expected actual type number, got %s", diffs[0].ActualType)
	}
}

func TestShapeMatch_ArrayElementShape(t *testing.T) {
	expected := json.RawMessage(`{"items": [{"id": "a", "value": 1}]}`)
	actual := json.RawMessage(`{"items": [{"id": "b", "value": 2}, {"id": "c", "value": 3}]}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_ArrayElementTypeMismatch(t *testing.T) {
	expected := json.RawMessage(`{"items": [{"id": "a", "count": 1}]}`)
	actual := json.RawMessage(`{"items": [{"id": "b", "count": "one"}]}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch, got match")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].Path != "$.items[0].count" {
		t.Errorf("expected path $.items[0].count, got %s", diffs[0].Path)
	}
}

func TestShapeMatch_EmptyObjects(t *testing.T) {
	expected := json.RawMessage(`{}`)
	actual := json.RawMessage(`{}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match for empty objects, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_EmptyObjectMatchesPopulated(t *testing.T) {
	// Empty expected matches any object (no required fields).
	expected := json.RawMessage(`{}`)
	actual := json.RawMessage(`{"id": "abc", "name": "test"}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match (empty expected), got diffs: %+v", diffs)
	}
}

func TestShapeMatch_NullHandling(t *testing.T) {
	// null matches null.
	expected := json.RawMessage(`{"value": null}`)
	actual := json.RawMessage(`{"value": null}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match for null values, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_NullVsString(t *testing.T) {
	expected := json.RawMessage(`{"value": "hello"}`)
	actual := json.RawMessage(`{"value": null}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch (string vs null), got match")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].ExpectedType != "string" || diffs[0].ActualType != "null" {
		t.Errorf("expected string/null mismatch, got %s/%s", diffs[0].ExpectedType, diffs[0].ActualType)
	}
}

func TestShapeMatch_EmptyArrays(t *testing.T) {
	expected := json.RawMessage(`{"items": []}`)
	actual := json.RawMessage(`{"items": []}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match for empty arrays, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_EmptyExpectedArrayMatchesPopulated(t *testing.T) {
	// Empty expected array matches any array.
	expected := json.RawMessage(`{"items": []}`)
	actual := json.RawMessage(`{"items": [{"id": "abc"}]}`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match (empty expected array), got diffs: %+v", diffs)
	}
}

func TestShapeMatch_TopLevelArray(t *testing.T) {
	expected := json.RawMessage(`[{"id": "a", "name": "test"}]`)
	actual := json.RawMessage(`[{"id": "b", "name": "other"}]`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match for top-level arrays, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_TopLevelPrimitive(t *testing.T) {
	expected := json.RawMessage(`"hello"`)
	actual := json.RawMessage(`"world"`)

	match, diffs := ShapeMatch(expected, actual)
	if !match {
		t.Errorf("expected match for string primitives, got diffs: %+v", diffs)
	}
}

func TestShapeMatch_MultipleDifferences(t *testing.T) {
	expected := json.RawMessage(`{"a": "str", "b": 42, "c": true}`)
	actual := json.RawMessage(`{"a": 1, "b": "wrong"}`)

	match, diffs := ShapeMatch(expected, actual)
	if match {
		t.Error("expected mismatch, got match")
	}
	// a: string vs number, b: number vs string, c: boolean vs missing
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d: %+v", len(diffs), diffs)
	}
}
