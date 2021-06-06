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

	paths, err := common.GetNotesPaths(folder, common.MdExtension)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found notes:", len(paths))

	for _, path := range paths {

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

		newPath := getFilepathOnlyId(path, id)
		err = os.Rename(string(path), newPath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s renamed to %s\n", path, newPath)
	}
}

func getFilepathOnlyId(oldPath common.Path, id string) string {
	return filepath.Dir(string(oldPath)) + "/" + id + ".md"
}

func parseNoteNameByPath(path common.Path) (isZettel bool, id, name string) {
	if filepath.Ext(string(path)) != ".md" {
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
