package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetNoteNameByNoteContent(content []string) (name string, err error) {
	for _, line := range content {
		titleYamlHeader := "title:"
		titleIdx := strings.Index(line, titleYamlHeader)

		hasYamlHeaderForTitle := titleIdx != -1
		hasFirstLevelHeader := strings.HasPrefix(line, "# ")

		//TODO it's better to rely on state machine and real parsing of yaml header
		if hasYamlHeaderForTitle {
			return line[titleIdx+len(titleYamlHeader)+1:], nil
		} else if hasFirstLevelHeader {
			return line[2:], nil
		}
	}

	return "", fmt.Errorf("not possible to detect and extract note name from file using yaml title or # header")
}

func IsZkId(id string) bool {
	//TODO customize zk id format/length/etc.
	if len(id) != 12 { //202005091607 = 4+2+2+2+2 = 12
		return false

	}

	_, err := strconv.Atoi(id)
	return err == nil
}

func ParseNoteFilename(filename string) (isZettel bool, id, name string) {
	spaceIndex := strings.Index(filename, " ")
	if spaceIndex == -1 {
		id = filename
	} else {
		id = filename[:spaceIndex]
	}

	if !IsZkId(id) {
		return false, "", ""
	}

	name = strings.TrimLeft(filename, id)
	name = strings.Trim(name, " ")
	return true, id, name
}

func GetYamlHeader(id string, name string) []string {
	return []string{
		"---",
		"title: " + strings.ToLower(name),
		"date: " + FormatIdAsIsoDate(id),
		"tags: ",
		"---",
	}
}

func FormatIdAsIsoDate(zkId string) string {
	date, err := time.Parse("200601021504", zkId)
	if err != nil {
		panic(err)
	}

	return date.Format("2006-01-02 15:04")
}
