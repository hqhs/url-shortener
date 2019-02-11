
.PHONY: all
all: build run

.PHONY: build
build:
	go build -o bin/service .

.PHONY: run
run:
	./bin/service

.PHONY: cache
cache:
	go build -i -o bin/service .

.PHONY: redis
redis:
	docker run -d -p 6379 redis
