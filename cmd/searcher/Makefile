build: container_build
	mv searcher.exe ../../bin/searcher.exe

container_build: generate
	rm -f *.sqlite
	rm -f *.sqlite.zstd
	go build --tags fts5 -buildvcs=false -o searcher.exe

test:
	go test

run: generate
	go run --tags fts5 *.go

check:
	staticcheck

generate:
	rm -f schema.sql
	cp ../../pkg/dyndb/schema.sql . 
	cp ../../pkg/dyndb/queries.sql . 
	cp ../../pkg/dyndb/sqlc.yaml . 
	sqlc generate
