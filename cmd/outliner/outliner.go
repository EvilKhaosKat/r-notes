package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"path/filepath"
	"strings"
	"time"
)

const (
	notesDelimiter    = "    "
	markdownLineBreak = "  "
	tag               = "#index"
)

func main() {
	path, folder, err := common.GetNoteFileArgument(common.MdExtension)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("generating outline for path", path)

	otherFiles, err := common.GetNotesPaths(folder, common.MdExtension)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(otherFiles))
	log.Println("parsing links")

	top := parseNoteHierarchy(path, otherFiles, 3)
	log.Printf("outline:\n")

	outline := getNotesOutline(top, "", nil)
	for _, line := range outline {
		fmt.Println(line)
	}

	resultId, resultPath := getResultPath(path)
	fmt.Printf("writing to %s\n", resultPath)

	indexTitle := fmt.Sprintf("index for '%s'", top.Name)

	resultContent := common.GetYamlHeader(resultId, indexTitle)
	resultContent = append(resultContent, getResultNoteHeader(resultPath, indexTitle), tag)
	resultContent = append(resultContent, outline...)

	common.WriteToFile(resultPath, resultContent)
}

func getResultNoteHeader(resultPath common.Path, title string) string {
	return fmt.Sprintf("# %s %s", common.GetFilename(resultPath), title)
}

func getResultPath(path common.Path) (id string, resultPath common.Path) {
	basePath := filepath.Dir(string(path))
	zkId := time.Now().Format("200601021504")
	return zkId, common.Path(
		fmt.Sprintf("%s/%s.md", basePath, zkId))
}

//TODO iterative version would be better, but lack of stdlib queue would decrease readability
func getNotesOutline(note *common.Note, padding string, result []string) []string {
	if note == nil {
		return result
	}

	noteLink := getNoteLink(note)
	result = append(result, fmt.Sprintf("%s- %s%s", padding, noteLink, markdownLineBreak))
	for _, child := range note.Links {
		result = getNotesOutline(child, padding+notesDelimiter, result)
	}

	return result
}

//TODO extract and reuse with logic from common pkg
func getNoteLink(note *common.Note) string {
	firstSpaceIndex := strings.Index(note.Name, " ")
	if firstSpaceIndex != -1 && common.IsZkId(note.Name[:firstSpaceIndex]) {
		return fmt.Sprintf("%s [[%s]]", note.Name[firstSpaceIndex+1:], note.Name[:firstSpaceIndex])
	} else {
		return fmt.Sprintf("%s [[%s]]", note.Name, note.Id)
	}
}

func parseNoteHierarchy(path common.Path, paths []common.Path, levelsLeft int) *common.Note {
	if levelsLeft == 0 {
		return nil
	}

	content, err := common.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	id, noteName, err := common.GetNoteNameByPath(path)
	if err != nil {
		panic(err)
	}

	note := common.NewNote(id, noteName, path, nil)

	linkedFiles := common.GetFilesByWikiLinks(path, paths, common.GetWikiLinks(content))
	for _, linkedFile := range linkedFiles {
		child := parseNoteHierarchy(linkedFile, paths, levelsLeft-1)
		if child != nil {
			note.Links = append(note.Links, child)
		}
	}

	return note
}
