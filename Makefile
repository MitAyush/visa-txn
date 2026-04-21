DOCKER_IMAGE := visa-txn
DOCKER_CONTAINER := visa-txn-app

MOCKERY_VERSION := v3.7.0

start: build run

deps:
	go mod tidy
	go mod vendor

clean:
	rm -f visa-txn
	rm -f storage/visa.db

build:
	go build -o visa-txn cmd/server/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -d --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE)

docker-clean:
	docker stop $(DOCKER_CONTAINER)
	docker rm $(DOCKER_CONTAINER)

run:
	./visa-txn

test:
	go test -v ./...

lint:
	go fmt ./...
	go vet ./...

mockery:
	go run github.com/vektra/mockery/v3@$(MOCKERY_VERSION)

.PHONY: start build run test lint clean deps docker-build docker-run docker-clean mockery