package main

import (
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"testing"
)

func Test_getPathWithIdAndTitle(t *testing.T) {
	type args struct {
		path common.Path
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generic note",
			args: struct {
				path common.Path
				name string
			}{
				path: "/ururu/202106062104.md",
				name: "generic note name"},
			want: "/ururu/202106062104 generic note name.md",
		},
		{
			name: "no tailing spaces",
			args: struct {
				path common.Path
				name string
			}{
				path: "/ururu/202106062104.md",
				name: "name "},
			want: "/ururu/202106062104 name.md",
		},
		{
			name: "no tailing .",
			args: struct {
				path common.Path
				name string
			}{
				path: "/ururu/202106062104.md",
				name: "name."},
			want: "/ururu/202106062104 name.md",
		},
		{
			name: "name has symbols which can't used in filename",
			args: struct {
				path common.Path
				name string
			}{
				path: "/ururu/202106062104.md",
				name: "separated/file\\name"},
			want: "/ururu/202106062104 separated file name.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPathWithIdAndTitle(tt.args.path, tt.args.name); got != tt.want {
				t.Errorf("getPathWithIdAndTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
