package code

import (
	"code/internal/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// GetDiff generates a formatted diff string comparing two configuration files.
// It accepts a slice of FileData containing file contents and their formats.
// The function parses each file according to its format (JSON or YAML),
// compares their key-value pairs, and returns a formatted string showing:
//   - Keys that were removed (prefixed with "- ")
//   - Keys that were added (prefixed with "+ ")
//   - Keys that were modified (shown as both removed and added)
//   - Keys that remain unchanged (prefixed with "  ")
//
// The output is sorted alphabetically by key names.
// Returns an error if file parsing fails.
func GetDiff(filesData []models.FileData) (string, error) {
	maps := make([]map[string]any, len(filesData))
	for i, fd := range filesData {
		maps[i] = make(map[string]any)
		if err := unmarshalFile(fd.Content, fd.Format, &maps[i]); err != nil {
			return "", err
		}
	}
	keys := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			keys[k] = struct{}{}
		}
	}
	sortedKey := make([]string, 0, len(keys))
	for k := range keys {
		sortedKey = append(sortedKey, k)
	}
	sort.Strings(sortedKey)

	old, new := maps[0], maps[1]
	var s strings.Builder
	s.WriteString("{\n")
	for _, key := range sortedKey {
		oldVal, inOld := old[key]
		newVal, inNew := new[key]

		switch {
		case inOld && !inNew:
			s.WriteString(printDiff("- ", key, oldVal))
		case !inOld && inNew:
			s.WriteString(printDiff("+ ", key, newVal))
		case fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal):
			s.WriteString(printDiff("- ", key, oldVal))
			s.WriteString(printDiff("+ ", key, newVal))
		default:
			s.WriteString(printDiff("  ", key, oldVal))
		}

	}
	s.WriteString("}")
	return s.String(), nil
}

func printDiff(sep, key string, val any) string {
	return fmt.Sprintf("  %s%s: %v\n", sep, key, val)
}

func unmarshalFile(data []byte, format string, v any) error {
	switch format {
	case ".json":
		if err := json.Unmarshal(data, v); err != nil {
			return err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format")
	}
	return nil
}
