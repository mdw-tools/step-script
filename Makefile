#!/usr/bin/make -f

VERSION := $(shell git describe)

test:
	go test -count=1 -short $(ARGS) ./...

install: test
	go install -ldflags="-X 'main.Version=$(VERSION)'"
