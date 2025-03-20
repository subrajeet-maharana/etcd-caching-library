# File: Makefile
.PHONY: build test clean run

build:
	go build -o bin/etcd-caching-library ./cmd/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

run: build
	./bin/etcd-caching-library

lint:
	golangci-lint run

deps:
	go mod tidy
	go mod download

docker-build:
	docker build -t etcd-caching-library .

docker-run:
	docker run -p 8080:8080 etcd-caching-library