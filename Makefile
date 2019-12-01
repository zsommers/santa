# Remove a bunch of old-school make cruft
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

SOURCE=$(shell find . -iname "*.go")

.PHONY: all
all: | build

build: $(SOURCE)
	go build -o $(GOBIN)/santa github.com/zsommers/santa

