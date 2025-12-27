package parsers

import (
	"code/internal/code"
	"code/internal/models"
	"fmt"
	"os"
	"strings"
)

// ParseByPaths reads JSON or YAML files from the given paths and generates a formatted
// diff showing the differences between them. It expects exactly two file paths.
// Supported file formats: .json, .yaml, .yml
// Files can be of different formats (e.g., comparing JSON with YAML is supported).
// It returns a string containing the diff output and an error if file reading or
// parsing fails.
func ParseByPaths(paths []string) (string, error) {
	filesData := make([]models.FileData, len(paths))
	for i, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		format, err := detectFormat(path)
		if err != nil {
			return "", err
		}
		filesData[i] = models.FileData{Content: data, Format: format}
	}
	out, err := code.GetDiff(filesData)
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
