# r-notes
## Auto-outliner
Generates an outline for a note.   
Relies on [[wiki-link]] format.   

For example with 3 levels:
![auto-outliner.png](auto-outliner.png)

### Build
` go build cmd/outliner/main.go`

### Usage
`outliner "path/to/note.md"`  
That will generate `/path/to/DATE Outline note.md`.