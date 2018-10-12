.PHONY: build-docker
build-docker:
	docker build -f Dockerfile-dev -t posix_mq .

.PHONY: build
build:
	go build example/exec/sender.go
	go build example/exec/receiver.go

.PHONY: run-sender
run-sender:
	go run example/exec/sender.go

.PHONY: run-receiver
run-receiver:
	go run example/exec/receiver.go
