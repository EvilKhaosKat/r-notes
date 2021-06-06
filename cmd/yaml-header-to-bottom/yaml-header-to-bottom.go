package main

import (
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
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

	for _, path := range paths {
		content, err := common.ReadFile(path)
		if err != nil {
			log.Printf("error while reading %s: %s\n", path, err)
			continue
		}

		updatedContent, canBeMoved := common.MoveHeaderFromTopToBottom(path, content)
		if canBeMoved {
			common.WriteToFile(path, updatedContent)
		}
	}
}
