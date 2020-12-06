# r-notes

## Outliner

Generates an outline for a markdown note with links.   
Relies on `[[wiki-link]]` format to find related notes.

For example with 3 levels depth (note -> links -> links of links)(The Archive app):
![outliner.png](outliner.png)

### Build

`go build ./cmd/outliner/`  
`go install ./cmd/outliner/`  
`go run ./cmd/outliner/`  

`go test ./cmd/outliner/`

### Usage

`outliner "path/to/note.md"`  
will generate `/path/to/DATE Outline note.md`.