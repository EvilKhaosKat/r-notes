package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	filename, folder, err := getFile()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("generating outline for file", filename)

	otherFiles, err := getMdFiles(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("found .md files:", len(otherFiles))
	log.Println("parsing links")

	content, err := readFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	links := getWikiLinks(content)
	fmt.Printf("%s\n", links)
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
	if !strings.HasSuffix(filename, ".md") {
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
