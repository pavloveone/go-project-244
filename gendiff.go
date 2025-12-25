package code

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ParseByPaths reads JSON files from the given paths and generates a formatted
// diff showing the differences between them. It expects exactly two file paths.
// It returns a string containing the diff output and an error if file reading or
// parsing fails.
func ParseByPaths(paths []string) (string, error) {
	files := make([][]byte, len(paths))
	for i, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		files[i] = data
	}
	out, err := genDiff(files)
	if err != nil {
		return "", err
	}
	return out, nil
}

func genDiff(files [][]byte) (string, error) {
	maps := make([]map[string]any, len(files))
	for i, file := range files {
		maps[i] = make(map[string]any)
		if err := json.Unmarshal(file, &maps[i]); err != nil {
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
