package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const MdExtension = ".md"

type Path string

func GetNotesPaths(path, extension string) ([]Path, error) {
	var paths []Path

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == extension {
				paths = append(paths, Path(path))
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

//TODO Trie would be much better
func GetFilesByWikiLinks(currentFile string, files []string, wikiLinks []string) []string {
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

func GetFilename(path Path) string {
	filename := filepath.Base(string(path))
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func GetNoteNameByPath(path Path) (id, name string, err error) {
	isZettel, id, name := ParseNoteFilename(GetFilename(path))
	if isZettel && len(name) != 0 {
		return id, name, nil
	}

	content, err := ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("reading note %v failed: %w", path, err)
	}

	name, err = GetNoteNameByNoteContent(content)
	if err != nil {
		return "", "", err
	}

	return id, name, nil
}

func ReadFile(path Path) ([]string, error) {
	file, err := os.Open(string(path))
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

func WriteToFile(path Path, content []string) {
	f, err := os.Create(string(path))
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
