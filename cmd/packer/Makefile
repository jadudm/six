build: container_build
	mv packer.exe ../../bin/packer.exe

container_build: generate
	rm -f *.sqlite
	rm -f *.sqlite.zstd
	go build --tags fts5 -buildvcs=false -o packer.exe

test:
	go test

run:
	go run --tags fts5 *.go

check:
	staticcheck

generate:
	rm -f schema.sql
	cd ../../pkg/dyndb ; sqlc generate
	cp ../../pkg/dyndb/schema.sql .