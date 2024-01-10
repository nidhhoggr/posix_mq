GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build

.PHONY: build

build: build_simple build_bidirectional build_duplex build_duplex_lag

.PHONY: build_simple
build_simple:
	go build -o bin/simple example/simple/simple.go

.PHONY: build_bidirectional
build_bidirectional:
	go build -o bin/bidirectional example/bidirectional/bidirectional.go

.PHONY: build_duplex
build_duplex:
	go build -o bin/duplex example/duplex/duplex.go

.PHONY: build_duplex_lag
build_duplex_lag:
	go build -o bin/duplex_lag example/duplex_lag/duplex_lag.go


.PHONY: test
test: test_simple test_bidirectional test_duplex test_duplex_lag

.PHONY: test_simple
test_simple:
	./bin/simple

.PHONY: test_bidrectional
test_bidirectional:
	./bin/bidirectional

.PHONY: test_duplex
test_duplex:
	./bin/duplex

.PHONY: test_duplex_lag
test_duplex_lag:
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
