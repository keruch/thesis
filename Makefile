#!/usr/bin/make -f

# Common variables
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
BUILDDIR ?= $(CURDIR)/build
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)

#? build: build the observer binary
build:
	go build -o $(BUILDDIR)/tssd ./cmd/tssd
.PHONY: build

#? install: install the observer binary
install:
	go install ./cmd/tssd
.PHONY: install

#? test: run tests
test:
	go test ./...
.PHONY: test

#? help: Get more info on make commands.
help:
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@cat $(MAKEFILE_LIST) | sed -n 's/^#?//p' | column -t -s ':' |  sort | sed -e 's/^/ /'
.PHONY: help