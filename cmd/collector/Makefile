build:
	go build --tags fts5 -o ../../bin/crawler.exe

test:
	go test

run:
	go run --tags fts5 main.go

# Migrations
up:
	dbmate -u "sqlite:_ignore/crawler.sqlite" up
down:
	dbmate -u "sqlite:_ignore/crawler.sqlite" down