package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"os"
	"sort"
)

func main() {
	folder, err := getNotesFolderArg()
	if err != nil {
		log.Fatal(err)
	}

	notes, err := common.GetMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for _, path := range notes {
		id, name, err := common.GetNoteNameByPath(path)
		if err != nil {
			panic(err)
		}

		if id != "" {
			result = append(result, fmt.Sprintf("[[%s]] %s", id, name))
		}
	}

	sort.Strings(result)

	for _, entry := range result {
		fmt.Println(entry)
	}
}

func getNotesFolderArg() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("specify notes folder")
	}

	return os.Args[1], nil
}
