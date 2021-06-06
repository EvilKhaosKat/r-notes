package common

type Note struct {
	Id   string
	Name string
	Path Path

	Parent   *Note
	LinkedTo []*Note
}

func (n Note) String() string {
	return n.Name
}

func NewNote(id, name string, path Path, parent *Note, linkedTo []*Note) *Note {
	return &Note{Id: id, Name: name, Path: path, Parent: parent, LinkedTo: linkedTo}
}
