clean:
	rm -f bin/queue.exe

build: clean
	go build --tags fts5 -o bin/queue.exe cmd/queue/main.go 