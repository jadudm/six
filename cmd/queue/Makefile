build:
	go build -o ../../bin/queue.exe

test:
	go test

run:
	go run main.go

# Migrations
up:
	dbmate -u "sqlite:db/queue.sqlite" up
down:
	dbmate -u "sqlite:db/queue.sqlite" down