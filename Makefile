.PHONY: clean

clean:
	rm -f bin/queue.exe

build: clean
	cd cmd/indexer ; make build
	cd cmd/queue-server ; make build
	# docker build -f -t com.jadud/six:latest .
