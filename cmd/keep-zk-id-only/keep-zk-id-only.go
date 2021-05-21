package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"os"
	"path/filepath"
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

		isZettel, id, name := parseNoteNameByPath(path)
		if !isZettel {
			fmt.Printf("%s is not a zettel\n", path)
			continue
		}

		if name == "" {
			fmt.Printf("%s is already in 'zk id only' format\n", path)
			continue
		}

		content, err := common.ReadFile(path)
		if err != nil {
			panic(err)
		}

		header := common.GetYamlHeader(id, name)

		common.WriteToFile(path, append(header, content...))

		newFilename := getFilepathOnlyId(path, id)
		err = os.Rename(path, newFilename)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s renamed to %s\n", path, newFilename)
	}
}

func getFilepathOnlyId(note string, id string) string {
	return filepath.Dir(note) + "/" + id + ".md"
}

func parseNoteNameByPath(path string) (isZettel bool, id, name string) {
	if filepath.Ext(path) != ".md" {
		return false, "", ""
	}

	return common.ParseNoteFilename(common.GetFilename(path))
}

func getNotesFolderArg() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("specify notes folder")
	}

	return os.Args[1], nil
}
