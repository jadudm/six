build: container_build
	mv indexer.exe ../../bin/indexer.exe

container_build:
	go build -o indexer.exe

test:
	go test

run:
	go run *.go

check:
	staticcheck

# Migrations
up:
	dbmate -u "sqlite:_ignore/crawler.sqlite" up
down:
	dbmate -u "sqlite:_ignore/crawler.sqlite" down