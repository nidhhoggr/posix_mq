GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

all: build

.PHONY: build

build: build_simplex build_bidirectional build_duplex build_duplex_lag

.PHONY: build_simplex
build_simplex:
	go build -o bin/simplex_sender example/simplex/sender/main.go
	go build -o bin/simplex_receiver example/simplex/receiver/main.go

.PHONY: build_bidirectional
build_bidirectional:
	go build -o bin/bidirectional_sender example/bidirectional/sender/main.go
	go build -o bin/bidirectional_responder example/bidirectional/responder/main.go

.PHONY: build_duplex
build_duplex:
	go build -o bin/duplex_sender example/duplex/sender/main.go
	go build -o bin/duplex_responder example/duplex/responder/main.go

.PHONY: build_duplex_lag
build_duplex_lag:
	go build -o bin/duplex_lag_sender example/duplex_lag/sender/main.go
	go build -o bin/duplex_lag_responder example/duplex_lag/responder/main.go


.PHONY: test
test: test_simplex test_bidirectional test_duplex test_duplex_lag

.PHONY: test_simplex
test_simplex:
	./bin/simplex_sender &
	sleep 1
	./bin/simplex_receiver

.PHONY: test_bidrectional
test_bidirectional:
	./bin/bidirectional_responder &
	sleep 1
	./bin/bidirectional_sender

.PHONY: test_duplex
test_duplex:
	./bin/duplex_responder &
	sleep 1
	./bin/duplex_sender

.PHONY: test_duplex_lag
test_duplex_lag:
	./bin/duplex_lag_responder &
	sleep 1
	./bin/duplex_lag_sender

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
