package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetMdFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".md" {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return files, nil
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

func GetFilename(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func GetNoteNameByPath(path string) (id, name string, err error) {
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

func ReadFile(path string) ([]string, error) {
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

func WriteToFile(path string, content []string) {
	f, err := os.Create(path)
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
