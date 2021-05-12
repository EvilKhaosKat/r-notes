package main

import (
	"fmt"
	"github.com/EvilKhaosKat/r-notes/pkg/common"
	"log"
	"os"
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

type Note struct {
	name string
	path string

	parent   *Note
	children []*Note
}

func (n Note) String() string {
	return n.name
}

func newNote(name string, path string, parent *Note, children []*Note) *Note {
	return &Note{name: name, path: path, parent: parent, children: children}
}

func main() {
	path, folder, err := getNoteFileArgument()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("generating outline for path", path)

	otherFiles, err := common.GetMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(otherFiles))
	log.Println("parsing links")

	top := parseNoteHierarchy(path, otherFiles, nil, 3)
	log.Printf("outline:\n")

	outline := printNotesOutline(top, "", nil)
	for _, line := range outline {
		fmt.Println(line)
	}

	resultFilename := getResultPath(path)
	fmt.Printf("writing to %s\n", resultFilename)

	resultContent := []string{getResultNoteHeader(resultFilename), tag}
	resultContent = append(resultContent, outline...)

	common.WriteToFile(resultFilename, resultContent)
}

func getResultNoteHeader(resultFilename string) string {
	return "# " + common.GetFilename(resultFilename)
}

func getResultPath(path string) string {
	basePath := filepath.Dir(path)
	return fmt.Sprintf("%s/%s.md",
		basePath,
		time.Now().Format("200601021504"))
}

//TODO iterative version would be better, but lack of stdlib queue would decrease readability
func printNotesOutline(note *Note, padding string, result []string) []string {
	if note == nil {
		return result
	}

	noteLink := getNoteLink(note)
	result = append(result, fmt.Sprintf("%s- %s%s", padding, noteLink, markdownLineBreak))
	for _, child := range note.children {
		result = printNotesOutline(child, padding+notesDelimiter, result)
	}

	return result
}

//TODO extract and reuse with logic from common pkg
func getNoteLink(note *Note) string {
	firstSpaceIndex := strings.Index(note.name, " ")
	if firstSpaceIndex != -1 && common.IsZkId(note.name[:firstSpaceIndex]) {
		return fmt.Sprintf("%s [[%s]]", note.name[firstSpaceIndex+1:], note.name[:firstSpaceIndex])
	} else {
		return note.name
	}
}

func parseNoteHierarchy(path string, files []string, parent *Note, levelsLeft int) *Note {
	if levelsLeft == 0 {
		return nil
	}

	content, err := common.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	noteName, err := common.GetNoteNameByPath(path)
	if err != nil {
		panic(err)
	}

	note := newNote(noteName, path, parent, nil)

	linkedFiles := common.GetFilesByWikiLinks(path, files, getWikiLinks(content))
	for _, linkedFile := range linkedFiles {
		child := parseNoteHierarchy(linkedFile, files, note, levelsLeft-1)
		if child != nil {
			note.children = append(note.children, child)
		}
	}

	return note
}

//getWikiLinks extracts [[LINK] from provided file content
func getWikiLinks(content []string) []string {
	set := make(map[string]struct{})          //lack of golang sets ;(
	re := regexp.MustCompile(`\[\[(.+?)\]\]`) //TODO compile once for app rather than once per file

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

func getNoteFileArgument() (string, string, error) {
	if len(os.Args) != 2 {
		return "", "", fmt.Errorf("specify path for generating outline")
	}

	filename := os.Args[1]
	if filepath.Ext(filename) == "md" {
		return "", "", fmt.Errorf("specify .md path for generating outline")
	}

	return filename, filepath.Dir(filename), nil
}
