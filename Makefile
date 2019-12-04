# Remove a bunch of old-school make cruft
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

SOURCE=$(shell find . -iname "*.go")

.PHONY: all
all: | build

build: $(SOURCE)
	go build -o $(GOBIN)/santa github.com/zsommers/santa/cmd/santa

.PHONY: lint
lint:
	golint ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: codeclimate
codeclimate:
	go test -coverprofile c.out ./...
	./cc-test-reporter format-coverage -t gocov --prefix github.com/zsommers/santa c.out
