GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build

.PHONY: build

build: build_simple build_bidirectional build_duplex build_duplex_lag

.PHONY: build_simple
build_simple:
	$(GO) build -o bin/simple example/simple/simple.go

.PHONY: build_bidirectional
build_bidirectional:
	$(GO) build -o bin/bidirectional example/bidirectional/bidirectional.go

.PHONY: build_duplex
build_duplex:
	$(GO) build -o bin/duplex example/duplex/duplex.go

.PHONY: build_duplex_lag
build_duplex_lag:
	$(GO) build -o bin/duplex_lag example/duplex_lag/duplex_lag.go

.PHONY: test
test: 
	$(GO) test -v

.PHONY: examples
examples: exa_simple exa_bidirectional exa_duplex exa_duplex_lag

.PHONY: exa_simple
exa_simple:
	./bin/simple

.PHONY: exa_bidrectional
exa_bidirectional:
	./bin/bidirectional

.PHONY: exa_duplex
exa_duplex:
	./bin/duplex

.PHONY: exa_duplex_lag
exa_duplex_lag:
	./bin/duplex_lag

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
