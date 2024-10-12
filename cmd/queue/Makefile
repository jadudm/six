build:
	go build -o ../../bin/queue.exe

test:
	go test handlers/*.go -v

run:
	go run *.go

check:
	staticcheck

# Migrations
# up:
# 	dbmate -u "sqlite:db/queue.sqlite" up
# down:
# 	dbmate -u "sqlite:db/queue.sqlite" down