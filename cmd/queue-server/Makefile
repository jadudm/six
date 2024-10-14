build: container_build
	cp queue-server.exe ../../bin/queue-server.exe

container_build:
	go build -buildvcs=false -o queue-server.exe

test:
	go test handlers/*.go -v

run:
	go run *.go

check:
	staticcheck
