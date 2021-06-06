package common

import (
	"fmt"
	"log"
	"strings"
)

const YamlDelimiter = "---"
const (
	yamlNotFound = iota
	readingYaml
)

type YamlHeader struct {
	Content  []string
	From, To int
}

func (y YamlHeader) String() string {
	return fmt.Sprintf("%d:%d\n%s", y.From, y.To, y.Content)
}

func NewYamlHeader(content []string, from int, to int) *YamlHeader {
	return &YamlHeader{Content: content, From: from, To: to}
}

func (y YamlHeader) Exists() bool {
	return y.From != -1
}

func ParseForYamlHeader(content []string) *YamlHeader {
	from := -1
	state := yamlNotFound

	for i, line := range content {
		if strings.HasPrefix(line, YamlDelimiter) {
			switch state {
			case yamlNotFound:
				state = readingYaml
				from = i
			case readingYaml:
				to := i
				return NewYamlHeader(content[from:to+1], from, to)
			}
		}
	}

	return NewYamlHeader(nil, -1, -1)
}

func MoveHeaderFromTopToBottom(path Path, content []string) ([]string, bool) {
	header := ParseForYamlHeader(content)
	if !header.Exists() {
		log.Printf("file %s doesn't have yaml header, skipping\n", path)
		return content, false
	}

	if header.From > 0 {
		log.Printf("file %s has yaml header not on top, skikpping\n", path)
		return content, false
	}

	body := content[header.To+1:]

	result := append(body, " ")
	return append(result, header.Content...), true
}
