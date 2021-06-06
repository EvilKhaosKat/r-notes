package common

import "regexp"

type Note struct {
	Id   string
	Name string
	Path Path

	//Backlinks *[]Note //TODO implement
	Links []*Note
}

func (n Note) String() string {
	return n.Name
}

func NewNote(id, name string, path Path, Links []*Note) *Note {
	return &Note{Id: id, Name: name, Path: path, Links: Links}
}

//GetWikiLinks extracts [[LINK]] from provided path content
//TODO make sure to guarantee order
func GetWikiLinks(content []string) []string {
	set := make(map[string]struct{})          //lack of golang sets ;(
	re := regexp.MustCompile(`\[\[(.+?)\]\]`) //TODO compile once for app rather than once per path

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
