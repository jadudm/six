.PHONY: clean

clean:
	rm -f bin/queue.exe

build: clean
	cd cmd/indexer ; make build
	cd cmd/queue-server ; make build
	# docker build -f -t com.jadud/six:latest .

make docker_build:
	docker build -f Dockerfile.builder --tag jadud/builder .
	docker build -f Dockerfile.indexer --tag jadud/indexer .
	docker build -f Dockerfile.queue-server --tag jadud/queue-server .