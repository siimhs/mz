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

.PHONY: build test vet tidy fmt go fakemz e2e e2e-update

build: $(CACHE)
	$(DOCKER_GO) go build -o mz .

test: $(CACHE)
	$(DOCKER_GO) go test ./...

# --- End-to-end tests (see e2e/README.md): the built binary against a fake MZ ---
# make e2e                        run all tests
# make e2e-update                 rewrite e2e/testdata/golden from actual output
E2E_GO = $(DOCKER_GO) sh -c

e2e: build $(CACHE)
	$(E2E_GO) 'cd e2e && go test -count=1 ./...'

e2e-update: build $(CACHE)
	$(E2E_GO) 'cd e2e && go test -count=1 . -update'

# Fixture tool for e2e/record.sh
fakemz: $(CACHE)
	$(E2E_GO) 'cd e2e && go build -o bin/fakemz ./cmd/fakemz'

vet: $(CACHE)
	$(DOCKER_GO) sh -c 'go vet ./... && cd e2e && go vet ./...'

tidy: $(CACHE)
	$(DOCKER_GO) go mod tidy

fmt: $(CACHE)
	$(DOCKER_GO) gofmt -l -w .

# Arbitrary go command: make go ARGS="list ./..."
go: $(CACHE)
	$(DOCKER_GO) go $(ARGS)

$(CACHE):
	mkdir -p $(CACHE)
