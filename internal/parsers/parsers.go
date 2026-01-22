package parsers

import (
	"code"
	"fmt"
)

// ParseByPaths reads JSON or YAML files from the given paths and generates a formatted
// diff showing the differences between them. It expects exactly two file paths and
// an output format string.
// Supported file formats: .json, .yaml, .yml
// Supported output formats: "stylish", "plain"
// Files can be of different formats (e.g., comparing JSON with YAML is supported).
// It returns a string containing the diff output and an error if file reading,
// parsing, or formatting fails.
func ParseByPaths(paths []string, format string) (string, error) {
	if len(paths) != 2 {
		return "", fmt.Errorf("expected exactly 2 paths, got %d", len(paths))
	}
	return code.GenDiff(paths[0], paths[1], format)
}
