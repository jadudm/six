build: container_build
	mv searcher.exe ../../bin/searcher.exe

container_build: generate
	go build --tags fts5 -buildvcs=false -o searcher.exe

test:
	go test

run: generate
	go run --tags fts5 *.go

check:
	staticcheck

generate:
	cd ../../pkg/dyndb ; sqlc generate
	cp ../../pkg/dyndb/schema.sql .