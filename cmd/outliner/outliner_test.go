package main

import (
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"reflect"
	"testing"
)

func Test_getFilesByLinks(t *testing.T) {
	type args struct {
		currentFile common.Path
		files       []common.Path
		wikiLinks   []string
	}
	tests := []struct {
		name string
		args args
		want []common.Path
	}{
		{
			name: "main",
			args: args{
				currentFile: "path.md",
				files:       []common.Path{"path.md", "first.md", "second.md", "third.md"},
				wikiLinks:   []string{"first", "third"}},
			want: []common.Path{"first.md", "third.md"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := common.GetFilesByWikiLinks(tt.args.currentFile, tt.args.files, tt.args.wikiLinks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilesByWikiLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNotesOutline(t *testing.T) {
	type args struct {
		note    *common.Note
		padding string
		result  []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "main",
			args: args{
				note: common.NewNote(
					"202105122138",
					"note",
					"/path/to/202105122138.md",
					[]*common.Note{
						{
							Id:   "202105122139",
							Name: "child",
							Path: "path/to/202105122139.md",
						},
					}),
				padding: "",
				result:  []string{},
			},
			want: []string{"- note [[202105122138]]  ", "    - child [[202105122139]]  "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNotesOutline(tt.args.note, tt.args.padding, tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNotesOutline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNoteName(t *testing.T) {
	type args struct {
		path common.Path
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "main",
			args: args{
				path: "/somewhere/path/note.md",
			},
			want: "note",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := common.GetFilename(tt.args.path); got != tt.want {
				t.Errorf("GetFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWikiLinks(t *testing.T) {
	type args struct {
		content []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				content: []string{"text [[link1]] and", "[[link2]] another link", "duplicate [[link1]] "},
			},
			want: []string{"link1", "link2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := common.GetWikiLinks(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWikiLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
