package formatters

import (
	"code/internal/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	indentSize   = 4
	markerOffset = 2
)

// FormatStylish formats a diff tree in the stylish format.
// It takes a slice of DiffNode representing the internal diff tree structure
// and returns a human-readable string with proper indentation and markers:
//   - Keys that were removed are prefixed with "- "
//   - Keys that were added are prefixed with "+ "
//   - Keys that were modified are shown as both removed and added
//   - Keys that remain unchanged are prefixed with "  "
//   - Nested structures are properly indented with 4 spaces per level
//
// The output uses consistent formatting with alphabetically sorted keys
// at each nesting level. Values are JSON-encoded to ensure proper representation
// of strings, numbers, booleans, and null values.
func FormatStylish(nodes []models.DiffNode) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	formatNodes(nodes, 1, &sb)
	sb.WriteString("}")
	return sb.String()
}

func formatNodes(nodes []models.DiffNode, depth int, sb *strings.Builder) {
	for _, node := range nodes {
		switch node.Type {
		case models.NodeTypeAdded:
			writeNode(sb, depth, "+ ", node.Key, node.NewValue)
		case models.NodeTypeRemoved:
			writeNode(sb, depth, "- ", node.Key, node.OldValue)
		case models.NodeTypeChanged:
			writeNode(sb, depth, "- ", node.Key, node.OldValue)
			writeNode(sb, depth, "+ ", node.Key, node.NewValue)
		case models.NodeTypeUnchanged:
			writeNode(sb, depth, "  ", node.Key, node.OldValue)
		case models.NodeTypeNested:
			writeNestedNode(sb, depth, node.Key, node.Children)
		}
	}
}

func writeNode(sb *strings.Builder, depth int, marker, key string, value any) {
	indent := strings.Repeat(" ", depth*indentSize-markerOffset)
	sb.WriteString(indent)
	sb.WriteString(marker)
	sb.WriteString(key)
	sb.WriteString(": ")
	sb.WriteString(formatValue(value, depth))
	sb.WriteString("\n")
}

func writeNestedNode(sb *strings.Builder, depth int, key string, children []models.DiffNode) {
	indent := strings.Repeat(" ", depth*indentSize-markerOffset)
	sb.WriteString(indent)
	sb.WriteString("  ")
	sb.WriteString(key)
	sb.WriteString(": {\n")
	formatNodes(children, depth+1, sb)
	sb.WriteString(indent)
	sb.WriteString("  }\n")
}

func formatValue(value any, depth int) string {
	if value == nil {
		return "null"
	}

	if m, ok := value.(map[string]any); ok {
		return formatMap(m, depth)
	}

	if s, ok := value.(string); ok {
		return s
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(bytes)
}

func formatMap(m map[string]any, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")

	keys := getSortedKeys(m)

	baseIndent := strings.Repeat(" ", (depth+1)*indentSize)
	for _, key := range keys {
		sb.WriteString(baseIndent)
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(formatValue(m[key], depth+1))
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat(" ", depth*indentSize))
	sb.WriteString("}")

	return sb.String()
}

func getSortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
