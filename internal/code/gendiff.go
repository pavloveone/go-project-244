package code

import (
	"code/internal/formatters"
	"code/internal/models"
	"encoding/json"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// GetDiff generates a formatted diff string comparing two configuration files.
// It accepts a slice of FileData containing file contents and their formats,
// and a format string specifying the output format.
// The function parses each file according to its format (JSON or YAML),
// compares their key-value pairs recursively, and returns a formatted string.
//
// Supported output formats:
//   - "stylish": Hierarchical format with indentation and markers
//   - "plain": Flat text format with property paths
//
// The output is sorted alphabetically by key names at each level.
// Returns an error if file parsing or formatting fails.
func GetDiff(filesData []models.FileData, format string) (string, error) {
	maps := make([]map[string]any, len(filesData))
	for i, fd := range filesData {
		maps[i] = make(map[string]any)
		if err := unmarshalFile(fd.Content, fd.Format, &maps[i]); err != nil {
			return "", err
		}
	}

	old, new := maps[0], maps[1]
	diffTree := BuildDiffTree(old, new)
	return formatters.Format(diffTree, format)
}

// BuildDiffTree recursively builds a diff tree comparing two maps.
// It analyzes the keys present in both maps and creates a slice of DiffNode
// representing the differences. The function handles:
//   - Added keys (present only in the new map)
//   - Removed keys (present only in the old map)
//   - Changed values (same key, different non-map values)
//   - Unchanged values (same key and value)
//   - Nested structures (same key with map values in both - recursively compared)
//
// Keys are sorted alphabetically at each level to ensure consistent output.
// Returns a slice of DiffNode representing the complete diff tree.
func BuildDiffTree(old, new map[string]any) []models.DiffNode {
	keys := make(map[string]struct{})
	for k := range old {
		keys[k] = struct{}{}
	}
	for k := range new {
		keys[k] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	nodes := make([]models.DiffNode, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		oldVal, inOld := old[key]
		newVal, inNew := new[key]

		node := models.DiffNode{Key: key}

		switch {
		case inOld && !inNew:
			node.Type = models.NodeTypeRemoved
			node.OldValue = oldVal

		case !inOld && inNew:
			node.Type = models.NodeTypeAdded
			node.NewValue = newVal

		case inOld && inNew:
			oldMap, oldIsMap := oldVal.(map[string]any)
			newMap, newIsMap := newVal.(map[string]any)

			if oldIsMap && newIsMap {
				node.Type = models.NodeTypeNested
				node.Children = BuildDiffTree(oldMap, newMap)
			} else if !valuesEqual(oldVal, newVal) {
				node.Type = models.NodeTypeChanged
				node.OldValue = oldVal
				node.NewValue = newVal
			} else {
				node.Type = models.NodeTypeUnchanged
				node.OldValue = oldVal
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func valuesEqual(a, b any) bool {
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)
	if aIsMap && bIsMap {
		if len(aMap) != len(bMap) {
			return false
		}
		for k, v := range aMap {
			bv, ok := bMap[k]
			if !ok || !valuesEqual(v, bv) {
				return false
			}
		}
		return true
	}

	return a == b
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
