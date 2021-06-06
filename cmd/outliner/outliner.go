package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	notesDelimiter    = "    "
	markdownLineBreak = "  "
	tag               = ":index:\n#index"
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

	top := parseNoteHierarchy(path, otherFiles, nil, 3)
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
	for _, child := range note.LinkedTo {
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

func parseNoteHierarchy(path common.Path, paths []common.Path, parent *common.Note, levelsLeft int) *common.Note {
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

	note := common.NewNote(id, noteName, path, parent, nil)

	linkedFiles := common.GetFilesByWikiLinks(path, paths, getWikiLinks(content))
	for _, linkedFile := range linkedFiles {
		child := parseNoteHierarchy(linkedFile, paths, note, levelsLeft-1)
		if child != nil {
			note.LinkedTo = append(note.LinkedTo, child)
		}
	}

	return note
}

//TODO make sure to guarantee order
//getWikiLinks extracts [[LINK] from provided path content
func getWikiLinks(content []string) []string {
	set := make(map[string]struct{})          //lack of golang sets ;(
	re := regexp.MustCompile(`\[\[(.+?)\]\]`) //TODO compile once for app rather than once per path

	for _, line := range content {
		for _, match := range re.FindAllStringSubmatch(line, -1) {
			link := match[1]
			set[link] = struct{}{}
		}
	}

	var links []string
	for link := range set {
		links = append(links, link)
	}

	return links
}
