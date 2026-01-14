package code

import (
	"code/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name    string
		files   []models.FileData
		want    string
		wantErr bool
	}{
		{
			name: "simple difference",
			files: []models.FileData{
				{Content: []byte(`{"a": 1, "b": 2}`), Format: ".json"},
				{Content: []byte(`{"a": 1, "c": 3}`), Format: ".json"},
			},
			want: "{\n    a: 1\n  - b: 2\n  + c: 3\n}",
		},
		{
			name: "value changes",
			files: []models.FileData{
				{Content: []byte(`{"key": "old"}`), Format: ".json"},
				{Content: []byte(`{"key": "new"}`), Format: ".json"},
			},
			want: "{\n  - key: \"old\"\n  + key: \"new\"\n}",
		},
		{
			name: "empty objects",
			files: []models.FileData{
				{Content: []byte(`{}`), Format: ".json"},
				{Content: []byte(`{}`), Format: ".json"},
			},
			want: "{\n}",
		},
		{
			name: "number type change",
			files: []models.FileData{
				{Content: []byte(`{"num": 42}`), Format: ".json"},
				{Content: []byte(`{"num": 43}`), Format: ".json"},
			},
			want: "{\n  - num: 42\n  + num: 43\n}",
		},
		{
			name: "boolean values",
			files: []models.FileData{
				{Content: []byte(`{"flag": true}`), Format: ".json"},
				{Content: []byte(`{"flag": false}`), Format: ".json"},
			},
			want: "{\n  - flag: true\n  + flag: false\n}",
		},
		{
			name: "invalid json in first file",
			files: []models.FileData{
				{Content: []byte(`{invalid}`), Format: ".json"},
				{Content: []byte(`{"a": 1}`), Format: ".json"},
			},
			wantErr: true,
		},
		{
			name: "invalid json in second file",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`not json`), Format: ".json"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := GenDiff(tt.files, "stylish")

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.want, got)
		})
	}
}

func TestPrintDiff(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		key  string
		val  any
		want string
	}{
		{
			name: "unchanged value",
			sep:  "  ",
			key:  "host",
			val:  "example.com",
			want: "    host: example.com\n",
		},
		{
			name: "added value",
			sep:  "+ ",
			key:  "port",
			val:  8080,
			want: "  + port: 8080\n",
		},
		{
			name: "removed value",
			sep:  "- ",
			key:  "debug",
			val:  true,
			want: "  - debug: true\n",
		},
		{
			name: "nil value",
			sep:  "  ",
			key:  "optional",
			val:  nil,
			want: "    optional: <nil>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got := printDiff(tt.sep, tt.key, tt.val)
			r.Equal(tt.want, got)
		})
	}
}

func TestGenDiffNested(t *testing.T) {
	tests := []struct {
		name    string
		files   []models.FileData
		want    string
		wantErr bool
	}{
		{
			name: "nested structures simple",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 2}}`), Format: ".json"},
			},
			want: `{
    a: {
      - b: 1
      + b: 2
    }
}`,
		},
		{
			name: "nested with added key",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 1, "c": 2}}`), Format: ".json"},
			},
			want: `{
    a: {
        b: 1
      + c: 2
    }
}`,
		},
		{
			name: "nested with removed key",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1, "c": 2}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
			},
			want: `{
    a: {
        b: 1
      - c: 2
    }
}`,
		},
		{
			name: "nested to non-nested",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": "string"}`), Format: ".json"},
			},
			want: `{
  - a: {
        b: 1
    }
  + a: "string"
}`,
		},
		{
			name: "non-nested to nested",
			files: []models.FileData{
				{Content: []byte(`{"a": "string"}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
			},
			want: `{
  - a: "string"
  + a: {
        b: 1
    }
}`,
		},
		{
			name: "deep nesting",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": {"c": 1}}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": {"c": 2}}}`), Format: ".json"},
			},
			want: `{
    a: {
        b: {
          - c: 1
          + c: 2
        }
    }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := GenDiff(tt.files, "stylish")

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.want, got)
		})
	}
}
