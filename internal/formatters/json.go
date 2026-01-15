package formatters

import (
	"code/internal/models"
	"encoding/json"
)

func FormatJSON(nodes []models.DiffNode) (string, error) {
	result := nodesToMap(nodes)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func nodesToMap(nodes []models.DiffNode) map[string]any {
	result := make(map[string]any)
	for _, node := range nodes {
		result[node.Key] = nodeToValue(node)
	}
	return result
}

func nodeToValue(node models.DiffNode) any {
	switch node.Type {
	case models.NodeTypeAdded:
		return map[string]any{
			"type":  "added",
			"value": node.NewValue,
		}
	case models.NodeTypeRemoved:
		return map[string]any{
			"type":  "removed",
			"value": node.OldValue,
		}
	case models.NodeTypeChanged:
		return map[string]any{
			"type":     "changed",
			"oldValue": node.OldValue,
			"newValue": node.NewValue,
		}
	case models.NodeTypeUnchanged:
		return map[string]any{
			"type":  "unchanged",
			"value": node.OldValue,
		}
	case models.NodeTypeNested:
		return map[string]any{
			"type":     "nested",
			"children": nodesToMap(node.Children),
		}
	}
	return nil
}
