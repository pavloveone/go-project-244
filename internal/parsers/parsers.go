package parsers

import (
	"code"
	"code/internal/models"
	"fmt"
	"os"
	"strings"
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
	filesData := make([]models.FileData, len(paths))
	for i, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		fileFormat, err := detectFormat(path)
		if err != nil {
			return "", err
		}
		filesData[i] = models.FileData{Content: data, Format: fileFormat}
	}
	out, err := code.GenDiff(filesData, format)
	if err != nil {
		return "", err
	}
	return out, nil
}

func detectFormat(path string) (string, error) {
	var formats = []string{".json", ".yaml", ".yml"}
	for _, f := range formats {
		if strings.HasSuffix(path, f) {
			return f, nil
		}
	}
	return "", fmt.Errorf("format has no support")
}
