build: container_build
	mv crawler.exe ../../bin/packer.exe

container_build:
	cd ../../pkg/dyndb ; sqlc generate
	cp ../../pkg/schema.sql .
	go build --tags fts5 -buildvcs=false -o packer.exe

test:
	go test

run:
	go run --tags fts5 *.go

check:
	staticcheck
