build:
	go build -o bin/cwlogsalert src/main.go

docker:
	docker build -t cwlogs-alert .

all: build docker
