package models

// NodeType represents the type of change in a diff node
type NodeType string

const (
	// NodeTypeAdded represents a key that was added
	NodeTypeAdded NodeType = "added"
	// NodeTypeRemoved represents a key that was removed
	NodeTypeRemoved NodeType = "removed"
	// NodeTypeChanged represents a key whose value changed
	NodeTypeChanged NodeType = "changed"
	// NodeTypeUnchanged represents a key whose value remained the same
	NodeTypeUnchanged NodeType = "unchanged"
	// NodeTypeNested represents a key whose value is a nested object in both files
	NodeTypeNested NodeType = "nested"
)

// DiffNode represents a single node in the diff tree
type DiffNode struct {
	Key      string      `json:"key"`
	Type     NodeType    `json:"type"`
	OldValue any         `json:"oldValue,omitempty"`
	NewValue any         `json:"newValue,omitempty"`
	Children []DiffNode  `json:"children,omitempty"`
}
