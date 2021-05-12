package common

import (
	"fmt"
	"strconv"
	"strings"
)

func GetNoteNameByNoteContent(content []string) (string, error) {
	for _, line := range content {
		titleYamlHeader := "title:"
		titleIdx := strings.Index(line, titleYamlHeader)

		hasYamlHeaderForTitle := titleIdx != -1
		hasFirstLevelHeader := strings.HasPrefix(line, "# ")

		//TODO it's better to rely on state machine and real parsing of yaml header
		if hasYamlHeaderForTitle {
			return line[titleIdx+len(titleYamlHeader):], nil
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
