build:
	go build --tags fts5 -o ../../bin/update-orchestrator.exe

test:
	go test

run:
	go run --tags fts5 main.go

check:
	staticcheck

# Migrations
up:
	dbmate -u "sqlite:_ignore/crawler.sqlite" up
down:
	dbmate -u "sqlite:_ignore/crawler.sqlite" down