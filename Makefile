# All Go tooling runs in Docker; only the resulting static binary runs on the host.
GO_IMAGE ?= golang:1.23
CACHE    := $(HOME)/.cache/mz-docker-go

DOCKER_GO = docker run --rm \
	-u $(shell id -u):$(shell id -g) \
	-v $(CURDIR):/src -w /src \
	-v $(CACHE):/cache \
	-e GOPATH=/cache/gopath -e GOCACHE=/cache/build \
	-e CGO_ENABLED=0 \
	$(GO_IMAGE)

.PHONY: build test vet tidy fmt go

build: $(CACHE)
	$(DOCKER_GO) go build -o mz .

test: $(CACHE)
	$(DOCKER_GO) go test ./...

vet: $(CACHE)
	$(DOCKER_GO) go vet ./...

tidy: $(CACHE)
	$(DOCKER_GO) go mod tidy

fmt: $(CACHE)
	$(DOCKER_GO) gofmt -l -w .

# Arbitrary go command: make go ARGS="list ./..."
go: $(CACHE)
	$(DOCKER_GO) go $(ARGS)

$(CACHE):
	mkdir -p $(CACHE)
