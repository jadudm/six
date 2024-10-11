clean:
	rm -f bin/queue.exe

build: clean
	cd cmd/queue && $(MAKE) build
