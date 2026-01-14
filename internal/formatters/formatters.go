package formatters

import (
	"code/internal/models"
	"fmt"
)

// Supported output formats
const (
	formatStylish = "stylish"
	formatPlain   = "plain"
	formatJson    = "json"
)

// Format formats a diff tree according to the specified format.
// It acts as a dispatcher, selecting the appropriate formatter based on the format parameter.
//
// Supported formats:
//   - "stylish": Hierarchical format with indentation and markers (default)
//   - "plain": Flat text format with property paths
//   - "json": json format
//
// Returns an error if an unknown format is specified.
func Format(nodes []models.DiffNode, format string) (string, error) {
	switch format {
	case formatStylish:
		return FormatStylish(nodes), nil
	case formatPlain:
		return FormatPlain(nodes), nil
	case formatJson:
		return FormatJSON(nodes)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
