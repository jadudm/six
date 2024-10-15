build: container_build
	mv packer.exe ../../bin/packer.exe

container_build:
	cd ../../pkg/dyndb ; sqlc generate
	cp ../../pkg/dyndb/schema.sql .
	go build --tags fts5 -buildvcs=false -o packer.exe

test:
	go test

run:
	go run --tags fts5 *.go

check:
	staticcheck
