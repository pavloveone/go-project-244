package formatters

import (
	"code/internal/models"
	"encoding/json"
)

// FormatJSON formats a diff tree in JSON format.
// It takes a slice of DiffNode representing the internal diff tree structure
// and returns a JSON string representation of the entire diff tree.
// The JSON output preserves the hierarchical structure of the diff,
// including node types, keys, values, and children.
//
// This format is useful for:
//   - Programmatic processing of diff results
//   - Integration with other tools and systems
//   - Storing diff results in a structured format
//   - Further analysis or transformation
//
// The output is pretty-printed with 2-space indentation for readability.
func FormatJSON(nodes []models.DiffNode) (string, error) {
	bytes, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
