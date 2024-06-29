# Project parameters
BINARY_NAME=load-balancer
SRC_DIR=./cmd/load-balancer
PKG_DIR=./internal/handler
TEST_DIRS=$(shell go list ./... | grep -v /vendor/)
DOCKER_COMPOSE_FILE=docker-compose.yml

# Go parameters
GO=go
GINKGO=ginkgo

.PHONY: all build run test clean tidy docker-up docker-down

all: clean build

build:
	$(GO) build -o $(BINARY_NAME) $(SRC_DIR)

run: build
	./$(BINARY_NAME)

test:
	$(GINKGO) $(TEST_DIRS)

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

tidy:
	$(GO) mod tidy

docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Convenience targets for common tasks
format:
	$(GO) fmt ./...

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

coverage:
	$(GO) test -coverprofile=coverage.out $(TEST_DIRS)
	$(GO) tool cover -html=coverage.out

ci: tidy vet format lint test coverage

# Running specific tests
test-specific:
	$(GINKGO) -focus="LoadBalancer RoundRobin" $(TEST_DIRS)
	$(GINKGO) -focus="LoadBalancer ServeHTTP" $(TEST_DIRS)

# Running docker-compose with logs
docker-logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f
