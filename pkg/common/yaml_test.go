package common

import (
	"reflect"
	"testing"
)

func TestParseForYamlHeader(t *testing.T) {
	type args struct {
		content []string
	}
	tests := []struct {
		name string
		args args
		want *YamlHeader
	}{
		{
			name: "top yaml header",
			args: struct{ content []string }{content: []string{
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
				"# header",
				"some value",
			}},
			want: &YamlHeader{
				Content: []string{
					"---",
					"title: ",
					"date: ",
					"tags: ",
					"---"},
				From: 0,
				To:   4,
			},
		},
		{
			name: "bottom yaml header",
			args: struct{ content []string }{content: []string{
				"# header",
				"some value",
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
			}},
			want: &YamlHeader{
				Content: []string{
					"---",
					"title: ",
					"date: ",
					"tags: ",
					"---"},
				From: 2,
				To:   6,
			},
		},
		{
			name: "no yaml header",
			args: struct{ content []string }{content: []string{
				"# header",
				"some value",
			}},
			want: &YamlHeader{
				Content: nil,
				From:    -1,
				To:      -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseForYamlHeader(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseForYamlHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYamlHeader_Exists(t *testing.T) {
	type fields struct {
		Content []string
		From    int
		To      int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "existing yaml header",
			fields: struct {
				Content []string
				From    int
				To      int
			}{
				Content: []string{}, From: 0, To: 10},
			want: true,
		},
		{
			name: "non-existing yaml header",
			fields: struct {
				Content []string
				From    int
				To      int
			}{
				Content: []string{}, From: -1, To: -1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := YamlHeader{
				Content: tt.fields.Content,
				From:    tt.fields.From,
				To:      tt.fields.To,
			}
			if got := y.Exists(); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoveHeaderFromTopToBottom(t *testing.T) {
	type args struct {
		path    Path
		content []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			name: "header moved from top to bottom",
			args: struct {
				path    Path
				content []string
			}{path: "/some/path/file.md", content: []string{
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
				"# header",
				"some value",
			}},
			want: []string{
				"# header",
				"some value",
				" ",
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
			},
			want1: true,
		},
		{
			name: "header not on top - not moved",
			args: struct {
				path    Path
				content []string
			}{path: "/some/path/file.md", content: []string{
				"# header",
				"some value",
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
			}},
			want: []string{
				"# header",
				"some value",
				"---",
				"title: ",
				"date: ",
				"tags: ",
				"---",
			},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MoveHeaderFromTopToBottom(tt.args.path, tt.args.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MoveHeaderFromTopToBottom() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MoveHeaderFromTopToBottom() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
