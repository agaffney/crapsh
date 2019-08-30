# Basic Makefile for crapsh
BINARY=crapsh

.PHONY: all clean

all: $(BINARY)

clean:
	rm -f $(BINARY)

$(BINARY): $(shell find -name '*.go')
	go build

.PHONY: run test

#TEST_CMD="foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789"
#TEST_CMD=foo baz\nabc \"123 456\" 'd\nef' 789
#TEST_CMD=echo foo bar && echo bar'baz baz'foo&&echo abc 123
TEST_CMD=FOO=bar BAR=baz echo foo >/dev/null bar 2>/dev/foo baz && echo abc 123 | grep foo &
run:
	go run crapsh.go -c "$(TEST_CMD)"

test:
	find -type f -name '*_test.go' | xargs -r dirname | sort -u | while read package; do \
		echo $$package; \
		go test $$package; \
	done
