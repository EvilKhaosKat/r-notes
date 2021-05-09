package main

import (
	"reflect"
	"testing"
)

func Test_getFilesByLinks(t *testing.T) {
	type args struct {
		currentFile string
		files       []string
		wikiLinks   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "main",
			args: args{
				currentFile: "file.md",
				files:       []string{"file.md", "first.md", "second.md", "third.md"},
				wikiLinks:   []string{"first", "third"}},
			want: []string{"first.md", "third.md"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFilesByWikiLinks(tt.args.currentFile, tt.args.files, tt.args.wikiLinks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilesByWikiLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNoteHierarchy(t *testing.T) {
	type args struct {
		note    *Note
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
				note: newNote(
					"note",
					"/path/to/note.md",
					nil,
					[]*Note{
						{
							name:     "child",
							filename: "path/to/child.md",
							parent:   nil, //TODO should have link to parent - create separate method for data prep
						},
					}),
				padding: "",
				result:  []string{},
			},
			want: []string{"[[note]]  ", "....[[child]]  "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printNotesOutline(tt.args.note, tt.args.padding, tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("printNotesOutline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNoteName(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "main",
			args: args{
				file: "/somewhere/file/note.md",
			},
			want: "note",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFullNoteName(tt.args.file); got != tt.want {
				t.Errorf("GetFullNoteName() = %v, want %v", got, tt.want)
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
			if got := getWikiLinks(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWikiLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
