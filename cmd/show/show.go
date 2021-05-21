package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"os"
)

func main() {
	folder, err := getNotesFolderArg()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("reading notes at " + folder)

	notes, err := common.GetMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(notes))

	for _, path := range notes {
		id, name, err := common.GetNoteNameByPath(path)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[[%s]] %s\n", id, name)
	}
}

func getNotesFolderArg() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("specify notes folder")
	}

	return os.Args[1], nil
}
