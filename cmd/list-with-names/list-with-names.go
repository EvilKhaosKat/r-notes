package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"sort"
)

func main() {
	folder, err := common.GetNotesFolderArg()
	if err != nil {
		log.Fatal(err)
	}

	paths, err := common.GetNotesPaths(folder, common.MdExtension)
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for _, path := range paths {
		id, name, err := common.GetNoteNameByPath(path)
		if err != nil {
			log.Printf("cant extract note name from file '%s'", path)
			continue
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
