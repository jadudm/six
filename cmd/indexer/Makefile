build: container_build
	mv scraper.exe ../../bin/scraper.exe

container_build:
	go build -buildvcs=false -o scraper.exe

test:
	go test

run:
	go run *.go

check:
	staticcheck
