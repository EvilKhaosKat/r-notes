package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	for _, noteFilename := range notes {

		isZettel, id, name := parseNoteNameByFilename(noteFilename)
		if !isZettel {
			fmt.Printf("%s is not a zettel\n", noteFilename)
			continue
		}

		if name == "" {
			fmt.Printf("%s is already in 'zk id only' format\n", noteFilename)
			continue
		}

		content, err := common.ReadFile(noteFilename)
		if err != nil {
			panic(err)
		}

		header := []string{
			"---",
			"title: " + strings.ToLower(name),
			"date: " + formatIdAsDate(id),
			"tags: ",
			"---",
		}

		common.WriteToFile(noteFilename, append(header, content...))

		newFilename := getFilepathOnlyId(noteFilename, id)
		err = os.Rename(noteFilename, newFilename)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s renamed to %s\n", noteFilename, newFilename)
	}
}

func getFilepathOnlyId(note string, id string) string {
	return filepath.Dir(note) + "/" + id + ".md"
}

func formatIdAsDate(zkId string) string {
	date, err := time.Parse("200601021504", zkId)
	if err != nil {
		panic(err)
	}

	return date.Format("2006-01-02 15:04")
}

func parseNoteNameByFilename(filename string) (isZettel bool, id, name string) {
	if filepath.Ext(filename) != ".md" {
		return false, "", ""
	}

	return common.ParseNoteFilename(common.GetFilename(filename))
}

func getNotesFolderArg() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("specify notes folder")
	}

	return os.Args[1], nil
}
