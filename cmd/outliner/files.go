package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
func getFilesByWikiLinks(currentFile string, files []string, wikiLinks []string) []string {
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

func GetFullNoteName(file string) string {
	fileName := filepath.Base(file)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getResultFilename(file string) string {
	basePath := filepath.Dir(file)
	return fmt.Sprintf("%s/%s %s %s.md",
		basePath,
		time.Now().Format("200601021504"),
		"Index",
		GetFullNoteName(file),
	)
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

func WriteToFile(filename string, content []string) {
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
