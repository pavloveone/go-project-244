package parsers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseByPaths(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		want    string
		wantErr bool
	}{
		{
			name:  "identical files",
			paths: []string{"../../testdata/fixture/file1.json", "../../testdata/fixture/file1.json"},
			want:  "{\n    follow: false\n    host: \"hexlet.io\"\n    proxy: \"123.234.53.22\"\n    timeout: 50\n}",
		},
		{
			name:  "different files with additions and deletions",
			paths: []string{"../../testdata/fixture/file1.json", "../../testdata/fixture/file2.json"},
			want:  "{\n  - follow: false\n    host: \"hexlet.io\"\n  - proxy: \"123.234.53.22\"\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}",
		},
		{
			name:  "empty files comparison",
			paths: []string{"../../testdata/fixture/empty.json", "../../testdata/fixture/empty.json"},
			want:  "{\n}",
		},
		{
			name:  "empty vs non-empty file",
			paths: []string{"../../testdata/fixture/empty.json", "../../testdata/fixture/single.json"},
			want:  "{\n  + key: \"value\"\n}",
		},
		{
			name:  "non-empty vs empty file",
			paths: []string{"../../testdata/fixture/single.json", "../../testdata/fixture/empty.json"},
			want:  "{\n  - key: \"value\"\n}",
		},
		{
			name:  "different data types",
			paths: []string{"../../testdata/fixture/types.json", "../../testdata/fixture/types.json"},
			want:  "{\n    boolean: true\n    null: null\n    number: 42\n    string: \"hello\"\n}",
		},
		{
			name:    "non-existent file",
			paths:   []string{"../../testdata/fixture/nonexistent.json", "../../testdata/fixture/file1.json"},
			wantErr: true,
		},
		{
			name:    "invalid json",
			paths:   []string{"../../testdata/fixture/invalid.json", "../../testdata/fixture/file1.json"},
			wantErr: true,
		},
		{
			name:  "identical yaml files",
			paths: []string{"../../testdata/fixture/file1.yaml", "../../testdata/fixture/file1.yaml"},
			want:  "{\n    follow: false\n    host: \"hexlet.io\"\n    proxy: \"123.234.53.22\"\n    timeout: 50\n}",
		},
		{
			name:  "different yaml files with additions and deletions",
			paths: []string{"../../testdata/fixture/file1.yaml", "../../testdata/fixture/file2.yaml"},
			want:  "{\n  - follow: false\n    host: \"hexlet.io\"\n  - proxy: \"123.234.53.22\"\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}",
		},
		{
			name:  "empty yaml files comparison",
			paths: []string{"../../testdata/fixture/empty.yaml", "../../testdata/fixture/empty.yaml"},
			want:  "{\n}",
		},
		{
			name:  "yaml vs json mixed format",
			paths: []string{"../../testdata/fixture/file1.yaml", "../../testdata/fixture/file2.json"},
			want:  "{\n  - follow: false\n    host: \"hexlet.io\"\n  - proxy: \"123.234.53.22\"\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}",
		},
		{
			name:  "json vs yaml mixed format",
			paths: []string{"../../testdata/fixture/file1.json", "../../testdata/fixture/file2.yaml"},
			want:  "{\n  - follow: false\n    host: \"hexlet.io\"\n  - proxy: \"123.234.53.22\"\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}",
		},
		{
			name:  "nested JSON files",
			paths: []string{"../../testdata/fixture/nested.json", "../../testdata/fixture/nested2.json"},
			want: `{
    common: {
      + follow: false
        setting1: "Value 1"
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: "blah blah"
      + setting5: {
            key5: "value5"
        }
        setting6: {
            doge: {
              - wow: ""
              + wow: "so much"
            }
            key: "value"
          + ops: "vops"
        }
    }
    group1: {
      - baz: "bas"
      + baz: "bars"
        foo: "bar"
      - nest: {
            key: "value"
        }
      + nest: "str"
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}`,
		},
		{
			name:  "nested YAML files",
			paths: []string{"../../testdata/fixture/nested.yaml", "../../testdata/fixture/nested2.yaml"},
			want: `{
    common: {
      + follow: false
        setting1: "Value 1"
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: "blah blah"
      + setting5: {
            key5: "value5"
        }
        setting6: {
            doge: {
              - wow: ""
              + wow: "so much"
            }
            key: "value"
          + ops: "vops"
        }
    }
    group1: {
      - baz: "bas"
      + baz: "bars"
        foo: "bar"
      - nest: {
            key: "value"
        }
      + nest: "str"
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := ParseByPaths(tt.paths, "stylish")

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.want, got)
		})
	}
}
