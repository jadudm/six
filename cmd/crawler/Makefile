build: container_build
	mv crawler.exe ../../bin/crawler.exe

container_build:
	go build -buildvcs=false -o crawler.exe

test:
	go test

run:
	go run *.go

check:
	staticcheck
