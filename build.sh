GOOS=linux GOARCH=amd64 go build -o outliner_linux github.com/EvilKhaosKat/r-notes/cmd/outliner && \
GOOS=darwin GOARCH=amd64 go build -o outliner_mac github.com/EvilKhaosKat/r-notes/cmd/outliner && \
GOOS=windows GOARCH=amd64 go build -o outliner_win github.com/EvilKhaosKat/r-notes/cmd/outliner