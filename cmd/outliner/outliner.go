package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	notesDelimiter    = "    "
	markdownLineBreak = "  "
	tag               = "#index"
)

type Note struct {
	name     string
	filename string

	parent   *Note
	children []*Note
}

func (n Note) String() string {
	return n.name
}

func newNote(name string, filename string, parent *Note, children []*Note) *Note {
	return &Note{name: name, filename: filename, parent: parent, children: children}
}

func main() {
	file, folder, err := getNoteFileArgument()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("generating outline for file", file)

	otherFiles, err := GetMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(otherFiles))
	log.Println("parsing links")

	top := parseNoteHierarchy(file, otherFiles, nil, 3)
	log.Printf("outline:\n")

	outline := printNotesOutline(top, "", nil)
	for _, line := range outline {
		fmt.Println(line)
	}

	resultFilename := getResultFilename(file)
	fmt.Printf("writing to %s\n", resultFilename)

	resultContent := []string{getResultNoteHeader(resultFilename), tag}
	resultContent = append(resultContent, outline...)

	WriteToFile(resultFilename, resultContent)
}

func getResultNoteHeader(resultFilename string) string {
	return "# " + GetFullNoteName(resultFilename)
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

func isId(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

func getNoteLink(note *Note) string {
	firstSpaceIndex := strings.Index(note.name, " ")
	if firstSpaceIndex != -1 && isId(note.name[:firstSpaceIndex]) {
		return fmt.Sprintf("%s [[%s]]", note.name[firstSpaceIndex+1:], note.name[:firstSpaceIndex])
	} else {
		return note.name
	}
}

func parseNoteHierarchy(file string, files []string, parent *Note, levelsLeft int) *Note {
	if levelsLeft == 0 {
		return nil
	}

	content, err := ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	note := newNote(GetFullNoteName(file), file, parent, nil)

	linkedFiles := getFilesByWikiLinks(file, files, getWikiLinks(content))
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
		return "", "", fmt.Errorf("specify filename for generating outline")
	}

	filename := os.Args[1]
	if filepath.Ext(filename) == "md" {
		return "", "", fmt.Errorf("specify .md filename for generating outline")
	}

	return filename, filepath.Dir(filename), nil
}
