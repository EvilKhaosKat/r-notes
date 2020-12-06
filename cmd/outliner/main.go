package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const notesDelimiter = "    "

type Note struct {
	name     string
	filename string

	parent   *Note
	children []*Note
}

func (n Note) String() string {
	return n.name
}

func newNote(name string, filename string, parent *Note) *Note {
	return &Note{name: name, filename: filename, parent: parent}
}

func main() {
	file, folder, err := getFile()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("generating outline for file", file)

	otherFiles, err := getMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(otherFiles))
	log.Println("parsing links")

	top := parseNoteHierarchy(file, otherFiles, nil, 3)
	log.Printf("outline:\n\n")

	hierarchy := getNoteHierarchy(top, "", nil)
	for _, line := range hierarchy {
		fmt.Println(line)
	}

	resultFilename := getResultFilename(file)
	fmt.Printf("writing to %s\n", resultFilename)

	resultContent := []string{"# " + top.name}
	resultContent = append(resultContent, hierarchy...)
	writeToFile(resultFilename, resultContent)
}

func getResultFilename(file string) string {
	basePath := filepath.Dir(file)
	return fmt.Sprintf("%s/%s %s %s.md",
		basePath,
		time.Now().Format("200601021504"),
		"Outline",
		getNoteName(file),
	)
}

func writeToFile(filename string, content []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range content {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func getNoteHierarchy(note *Note, delimiter string, result []string) []string {
	if note == nil {
		return result
	}

	result = append(result, fmt.Sprintf("%s[[%s]]", delimiter, note.String()))
	for _, child := range note.children {
		result = getNoteHierarchy(child, delimiter+notesDelimiter, result)
	}

	return result
}

func parseNoteHierarchy(file string, files []string, parent *Note, levelsLeft int) *Note {
	if levelsLeft == 0 {
		return nil
	}

	content, err := readFile(file)
	if err != nil {
		log.Fatal(err)
	}

	note := newNote(getNoteName(file), file, parent)

	linkedFiles := getFilesByLinks(file, files, getWikiLinks(content))
	for _, linkedFile := range linkedFiles {
		child := parseNoteHierarchy(linkedFile, files, note, levelsLeft-1)
		if child != nil {
			note.children = append(note.children, child)
		}
	}

	return note
}

func getNoteName(file string) string {
	fileName := filepath.Base(file)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getFilesByLinks(currentFile string, files []string, wikiLinks []string) []string {
	var linkedFiles []string

	for _, file := range files {
		for _, link := range wikiLinks {
			if file != currentFile && strings.Contains(file, link) {
				linkedFiles = append(linkedFiles, file)
			}
		}
	}

	return linkedFiles
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

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func getFile() (string, string, error) {
	if len(os.Args) != 2 {
		return "", "", fmt.Errorf("specify filename for generating outline")
	}

	filename := os.Args[1]
	if filepath.Ext(filename) == "md" {
		return "", "", fmt.Errorf("specify .md filename for generating outline")
	}

	return filename, filepath.Dir(filename), nil
}

func getMdFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return files, nil
}
