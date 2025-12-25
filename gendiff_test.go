package code

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
			paths: []string{"./testdata/fixture/file1.json", "./testdata/fixture/file1.json"},
			want:  "{\n    follow: false\n    host: hexlet.io\n    proxy: 123.234.53.22\n    timeout: 50\n}",
		},
		{
			name:  "different files with additions and deletions",
			paths: []string{"./testdata/fixture/file1.json", "./testdata/fixture/file2.json"},
			want:  "{\n  - follow: false\n    host: hexlet.io\n  - proxy: 123.234.53.22\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}",
		},
		{
			name:  "empty files comparison",
			paths: []string{"./testdata/fixture/empty.json", "./testdata/fixture/empty.json"},
			want:  "{\n}",
		},
		{
			name:  "empty vs non-empty file",
			paths: []string{"./testdata/fixture/empty.json", "./testdata/fixture/single.json"},
			want:  "{\n  + key: value\n}",
		},
		{
			name:  "non-empty vs empty file",
			paths: []string{"./testdata/fixture/single.json", "./testdata/fixture/empty.json"},
			want:  "{\n  - key: value\n}",
		},
		{
			name:  "different data types",
			paths: []string{"./testdata/fixture/types.json", "./testdata/fixture/types.json"},
			want:  "{\n    boolean: true\n    null: <nil>\n    number: 42\n    string: hello\n}",
		},
		{
			name:    "non-existent file",
			paths:   []string{"./testdata/fixture/nonexistent.json", "./testdata/fixture/file1.json"},
			wantErr: true,
		},
		{
			name:    "invalid json",
			paths:   []string{"./testdata/fixture/invalid.json", "./testdata/fixture/file1.json"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := ParseByPaths(tt.paths)

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.want, got)
		})
	}
}

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name    string
		files   [][]byte
		want    string
		wantErr bool
	}{
		{
			name: "simple difference",
			files: [][]byte{
				[]byte(`{"a": 1, "b": 2}`),
				[]byte(`{"a": 1, "c": 3}`),
			},
			want: "{\n    a: 1\n  - b: 2\n  + c: 3\n}",
		},
		{
			name: "value changes",
			files: [][]byte{
				[]byte(`{"key": "old"}`),
				[]byte(`{"key": "new"}`),
			},
			want: "{\n  - key: old\n  + key: new\n}",
		},
		{
			name: "empty objects",
			files: [][]byte{
				[]byte(`{}`),
				[]byte(`{}`),
			},
			want: "{\n}",
		},
		{
			name: "number type change",
			files: [][]byte{
				[]byte(`{"num": 42}`),
				[]byte(`{"num": 43}`),
			},
			want: "{\n  - num: 42\n  + num: 43\n}",
		},
		{
			name: "boolean values",
			files: [][]byte{
				[]byte(`{"flag": true}`),
				[]byte(`{"flag": false}`),
			},
			want: "{\n  - flag: true\n  + flag: false\n}",
		},
		{
			name: "invalid json in first file",
			files: [][]byte{
				[]byte(`{invalid}`),
				[]byte(`{"a": 1}`),
			},
			wantErr: true,
		},
		{
			name: "invalid json in second file",
			files: [][]byte{
				[]byte(`{"a": 1}`),
				[]byte(`not json`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := genDiff(tt.files)

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
