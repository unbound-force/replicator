// Package parity provides shape comparison tools for verifying that the Go
// rewrite of cyborg-swarm produces responses with the same JSON structure
// as the TypeScript original.
//
// Shape comparison checks field names and types (string, number, boolean,
// null, object, array) without comparing values. Extra fields in the actual
// response are acceptable (compatible superset), but missing fields are
// reported as mismatches.
package parity

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Difference describes a structural mismatch between expected and actual JSON.
type Difference struct {
	Path         string // JSONPath notation, e.g., "$.content[0].text"
	ExpectedType string // "string", "number", "boolean", "null", "object", "array"
	ActualType   string // same set, or "missing" if the key is absent
}

// ShapeMatch performs a recursive JSON shape comparison between expected and
// actual responses. It returns true if the actual response is a compatible
// superset of the expected shape (all expected fields present with matching
// types). Extra fields in actual are allowed and do not cause a mismatch.
//
// For arrays, only the first element's shape is compared (if non-empty).
// For primitives, only the type is compared (string/number/boolean/null).
func ShapeMatch(expected, actual json.RawMessage) (bool, []Difference) {
	var diffs []Difference
	compareShape("$", expected, actual, &diffs)
	return len(diffs) == 0, diffs
}

// jsonType returns the JSON type name for a raw value.
// Uses byte inspection for unambiguous types (null, boolean, string, array)
// and falls back to Unmarshal only for number vs object disambiguation.
func jsonType(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "null"
	}

	// Check for literal null first -- before any Unmarshal attempts,
	// because json.Unmarshal(null, &map) succeeds with a nil map.
	trimmed := string(raw)
	if trimmed == "null" {
		return "null"
	}

	// Use first byte to disambiguate unambiguous types.
	switch raw[0] {
	case '"':
		return "string"
	case '[':
		return "array"
	case '{':
		return "object"
	case 't', 'f':
		return "boolean"
	}

	// Must be a number (digits, minus sign, etc.).
	var n float64
	if json.Unmarshal(raw, &n) == nil {
		return "number"
	}

	return "unknown"
}

// compareShape recursively compares JSON shapes, accumulating differences.
func compareShape(path string, expected, actual json.RawMessage, diffs *[]Difference) {
	expType := jsonType(expected)
	actType := jsonType(actual)

	// Type mismatch at this level.
	if expType != actType {
		*diffs = append(*diffs, Difference{
			Path:         path,
			ExpectedType: expType,
			ActualType:   actType,
		})
		return
	}

	switch expType {
	case "object":
		compareObjects(path, expected, actual, diffs)
	case "array":
		compareArrays(path, expected, actual, diffs)
		// Primitives (string, number, boolean, null): type match is sufficient.
	}
}

// compareObjects compares two JSON objects by key set and recurses on values.
// Missing keys in actual are reported. Extra keys in actual are allowed.
func compareObjects(path string, expected, actual json.RawMessage, diffs *[]Difference) {
	var expObj, actObj map[string]json.RawMessage
	if err := json.Unmarshal(expected, &expObj); err != nil {
		return
	}
	if err := json.Unmarshal(actual, &actObj); err != nil {
		return
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(expObj))
	for k := range expObj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		childPath := fmt.Sprintf("%s.%s", path, key)
		expVal := expObj[key]

		actVal, exists := actObj[key]
		if !exists {
			*diffs = append(*diffs, Difference{
				Path:         childPath,
				ExpectedType: jsonType(expVal),
				ActualType:   "missing",
			})
			continue
		}

		compareShape(childPath, expVal, actVal, diffs)
	}

	// Extra keys in actual are OK -- compatible superset.
	// We intentionally do NOT report them as differences.
}

// compareArrays compares two JSON arrays by comparing the shape of their
// first elements (if both are non-empty). Empty arrays match any array.
func compareArrays(path string, expected, actual json.RawMessage, diffs *[]Difference) {
	var expArr, actArr []json.RawMessage
	if err := json.Unmarshal(expected, &expArr); err != nil {
		return
	}
	if err := json.Unmarshal(actual, &actArr); err != nil {
		return
	}

	// If expected array is empty, any array matches.
	if len(expArr) == 0 {
		return
	}

	// If expected has elements but actual is empty, that's still a shape match
	// for the array itself -- we can't compare element shapes.
	if len(actArr) == 0 {
		return
	}

	// Compare first element shapes.
	elemPath := fmt.Sprintf("%s[0]", path)
	compareShape(elemPath, expArr[0], actArr[0], diffs)
}
