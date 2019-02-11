
.PHONY: all
all: build

.PHONY: run
run: build bootstrap

.PHONY: build
build:
	go build -o bin/service .

.PHONY: bootstrap
bootstrap:
	./bin/service

.PHONY: cache
cache:
	go build -i -o bin/service .

.PHONY: redis
redis:
	docker run -p "6379:6379" redis
