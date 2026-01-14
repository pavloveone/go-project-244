package formatters

import (
	"code/internal/models"
	"fmt"
	"strings"
)

// FormatPlain formats a diff tree in the plain text format.
// It takes a slice of DiffNode representing the internal diff tree structure
// and returns a human-readable text output showing property changes:
//   - Added properties: "Property 'path' was added with value: X"
//   - Removed properties: "Property 'path' was removed"
//   - Changed properties: "Property 'path' was updated. From X to Y"
//   - Unchanged properties are not shown
//   - Nested objects show full path separated by dots (e.g., 'common.setting6.ops')
//   - Complex values (objects) are shown as [complex value]
//   - String values are wrapped in single quotes
//
// The output is sorted alphabetically by property path.
func FormatPlain(nodes []models.DiffNode) string {
	lines := formatPlainNodes(nodes, "")
	return strings.Join(lines, "\n")
}

func formatPlainNodes(nodes []models.DiffNode, parentPath string) []string {
	var lines []string

	for _, node := range nodes {
		path := buildPath(parentPath, node.Key)

		switch node.Type {
		case models.NodeTypeAdded:
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s", path, formatPlainValue(node.NewValue)))

		case models.NodeTypeRemoved:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", path))

		case models.NodeTypeChanged:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				path, formatPlainValue(node.OldValue), formatPlainValue(node.NewValue)))

		case models.NodeTypeNested:
			childLines := formatPlainNodes(node.Children, path)
			lines = append(lines, childLines...)

		}
	}

	return lines
}

func buildPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}
	return parentPath + "." + key
}

func isComplexValue(value any) bool {
	_, ok := value.(map[string]any)
	return ok
}

func formatPlainValue(value any) string {
	if isComplexValue(value) {
		return "[complex value]"
	}

	if value == nil {
		return "null"
	}

	if str, ok := value.(string); ok {
		return "'" + str + "'"
	}

	// For numbers and booleans
	return fmt.Sprintf("%v", value)
}
