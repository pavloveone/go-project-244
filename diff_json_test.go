package code

import (
	"code/internal/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiffJSON(t *testing.T) {
	tests := []struct {
		name    string
		files   []models.FileData
		want    string
		wantErr bool
	}{
		{
			name: "simple property added",
			files: []models.FileData{
				{Content: []byte(`{}`), Format: ".json"},
				{Content: []byte(`{"a": 1}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "added",
    "newValue": 1
  }
]`,
		},
		{
			name: "simple property removed",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`{}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "removed",
    "oldValue": 1
  }
]`,
		},
		{
			name: "property unchanged",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`{"a": 1}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "unchanged",
    "oldValue": 1
  }
]`,
		},
		{
			name: "property changed",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`{"a": 2}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "changed",
    "oldValue": 1,
    "newValue": 2
  }
]`,
		},
		{
			name: "nested structure",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 2}}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "nested",
    "children": [
      {
        "key": "b",
        "type": "changed",
        "oldValue": 1,
        "newValue": 2
      }
    ]
  }
]`,
		},
		{
			name: "multiple changes",
			files: []models.FileData{
				{Content: []byte(`{"a": 1, "b": 2, "c": 3}`), Format: ".json"},
				{Content: []byte(`{"a": 1, "b": 20, "d": 4}`), Format: ".json"},
			},
			want: `[
  {
    "key": "a",
    "type": "unchanged",
    "oldValue": 1
  },
  {
    "key": "b",
    "type": "changed",
    "oldValue": 2,
    "newValue": 20
  },
  {
    "key": "c",
    "type": "removed",
    "oldValue": 3
  },
  {
    "key": "d",
    "type": "added",
    "newValue": 4
  }
]`,
		},
		{
			name: "different value types",
			files: []models.FileData{
				{Content: []byte(`{"str": "hello", "num": 42, "bool": true, "null": null}`), Format: ".json"},
				{Content: []byte(`{"str": "world", "num": 43, "bool": false, "null": "value"}`), Format: ".json"},
			},
			want: `[
  {
    "key": "bool",
    "type": "changed",
    "oldValue": true,
    "newValue": false
  },
  {
    "key": "null",
    "type": "changed",
    "newValue": "value"
  },
  {
    "key": "num",
    "type": "changed",
    "oldValue": 42,
    "newValue": 43
  },
  {
    "key": "str",
    "type": "changed",
    "oldValue": "hello",
    "newValue": "world"
  }
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := GenDiff(tt.files, "json")

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)

			var gotJSON, wantJSON any
			r.NoError(json.Unmarshal([]byte(got), &gotJSON), "got should be valid JSON")
			r.NoError(json.Unmarshal([]byte(tt.want), &wantJSON), "want should be valid JSON")

			r.Equal(wantJSON, gotJSON)
		})
	}
}

func TestJSONFormatterValidOutput(t *testing.T) {
	files := []models.FileData{
		{Content: []byte(`{"common":{"setting2":200,"setting3":true,"setting6":{"key":"value","doge":{"wow":""}}},"group1":{"baz":"bas","foo":"bar","nest":{"key":"value"}},"group2":{"abc":12345,"deep":{"id":45}}}`), Format: ".json"},
		{Content: []byte(`{"common":{"follow":false,"setting1":"Value 1","setting3":null,"setting4":"blah blah","setting5":{"key5":"value5"},"setting6":{"key":"value","ops":"vops","doge":{"wow":"so much"}}},"group1":{"foo":"bar","baz":"bars","nest":"str"},"group3":{"deep":{"id":{"number":45}},"fee":100500}}`), Format: ".json"},
	}

	got, err := GenDiff(files, "json")
	require.NoError(t, err)

	var result any
	err = json.Unmarshal([]byte(got), &result)
	require.NoError(t, err, "Output should be valid JSON")

	arr, ok := result.([]any)
	require.True(t, ok, "Root element should be an array")
	require.NotEmpty(t, arr, "Array should not be empty")
}
