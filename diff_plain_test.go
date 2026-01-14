package code

import (
	"code/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiffPlain(t *testing.T) {
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
			want: "Property 'a' was added with value: 1",
		},
		{
			name: "simple property removed",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`{}`), Format: ".json"},
			},
			want: "Property 'a' was removed",
		},
		{
			name: "simple property updated",
			files: []models.FileData{
				{Content: []byte(`{"a": 1}`), Format: ".json"},
				{Content: []byte(`{"a": 2}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From 1 to 2",
		},
		{
			name: "string values with quotes",
			files: []models.FileData{
				{Content: []byte(`{"a": "old"}`), Format: ".json"},
				{Content: []byte(`{"a": "new"}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From 'old' to 'new'",
		},
		{
			name: "boolean values",
			files: []models.FileData{
				{Content: []byte(`{"a": true}`), Format: ".json"},
				{Content: []byte(`{"a": false}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From true to false",
		},
		{
			name: "null values",
			files: []models.FileData{
				{Content: []byte(`{"a": "value"}`), Format: ".json"},
				{Content: []byte(`{"a": null}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From 'value' to null",
		},
		{
			name: "nested path",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 2}}`), Format: ".json"},
			},
			want: "Property 'a.b' was updated. From 1 to 2",
		},
		{
			name: "deep nested path",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": {"c": 1}}}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": {"c": 2}}}`), Format: ".json"},
			},
			want: "Property 'a.b.c' was updated. From 1 to 2",
		},
		{
			name: "complex value added",
			files: []models.FileData{
				{Content: []byte(`{}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
			},
			want: "Property 'a' was added with value: [complex value]",
		},
		{
			name: "complex value removed",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{}`), Format: ".json"},
			},
			want: "Property 'a' was removed",
		},
		{
			name: "from complex to simple",
			files: []models.FileData{
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
				{Content: []byte(`{"a": "string"}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From [complex value] to 'string'",
		},
		{
			name: "from simple to complex",
			files: []models.FileData{
				{Content: []byte(`{"a": "string"}`), Format: ".json"},
				{Content: []byte(`{"a": {"b": 1}}`), Format: ".json"},
			},
			want: "Property 'a' was updated. From 'string' to [complex value]",
		},
		{
			name: "multiple changes",
			files: []models.FileData{
				{Content: []byte(`{"a": 1, "b": 2}`), Format: ".json"},
				{Content: []byte(`{"a": 3, "c": 4}`), Format: ".json"},
			},
			want: `Property 'a' was updated. From 1 to 3
Property 'b' was removed
Property 'c' was added with value: 4`,
		},
		{
			name: "nested complex file",
			files: []models.FileData{
				{Content: []byte(`{"common":{"setting2":200,"setting3":true,"setting6":{"key":"value","doge":{"wow":""}}},"group1":{"baz":"bas","foo":"bar","nest":{"key":"value"}},"group2":{"abc":12345,"deep":{"id":45}}}`), Format: ".json"},
				{Content: []byte(`{"common":{"follow":false,"setting1":"Value 1","setting3":null,"setting4":"blah blah","setting5":{"key5":"value5"},"setting6":{"key":"value","ops":"vops","doge":{"wow":"so much"}}},"group1":{"foo":"bar","baz":"bars","nest":"str"},"group3":{"deep":{"id":{"number":45}},"fee":100500}}`), Format: ".json"},
			},
			want: `Property 'common.follow' was added with value: false
Property 'common.setting1' was added with value: 'Value 1'
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := GenDiffFromData(tt.files, "plain")

			if tt.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.want, got)
		})
	}
}
