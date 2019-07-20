# Basic Makefile for crapsh
BINARY=crapsh

.PHONY:all

all: $(BINARY)

$(BINARY): $(shell find -name '*.go')
	go build

.PHONY: run test

run:
	go run crapsh.go

test:
	find -type f -name '*_test.go' | xargs -r dirname | sort -u | while read package; do \
		echo $$package; \
		go test $$package; \
	done
