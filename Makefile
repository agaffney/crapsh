# Basic Makefile for crapsh
BINARY=crapsh

all: $(BINARY)

$(BINARY):
	go build

test:
	find -type f -name '*_test.go' | xargs -r dirname | sort -u | while read package; do \
		echo $$package; \
		go test $$package; \
	done
