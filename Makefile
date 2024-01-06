GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build test

.PHONY: docker
docker:
	docker build -f Dockerfile-alpine -t posix_mq_alpine .

.PHONY: build

build: build_simplex build_duplex

.PHONY: build_simplex
build_simplex:
	go build -o bin/simplex_sender example/simplex/sender/main.go
	go build -o bin/simplex_receiver example/simplex/receiver/main.go

.PHONY: build_duplex
build_duplex:
	go build -o bin/duplex_sender example/duplex/sender/main.go
	go build -o bin/duplex_responder example/duplex/responder/main.go

.PHONY: test
test: test_simplex test_duplex

.PHONY: test_simplex
test_simplex:
	./bin/simplex_sender &
	./bin/simplex_receiver

.PHONY: test_duplex
test_duplex:
	./bin/duplex_responder &
	./bin/duplex_sender

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
  if [ -n "$$diff" ]; then \
    echo "Please run 'make fmt' and commit the result:"; \
    echo "$${diff}"; \
    exit 1; \
  fi;
