GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build

.PHONY: build

build: build_simple build_bidirectional 

.PHONY: build_simple
build_simple:
	$(GO) build -o bin/simple example/simple/simple.go

.PHONY: build_bidirectional
build_bidirectional:
	$(GO) build -o bin/bidirectional example/bidirectional/bidirectional.go

.PHONY: test
test: 
	$(GO) test -v

.PHONY: examples
examples: 
	./bin/simple
	./bin/bidirectional

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
