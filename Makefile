NAME     := wsbridge
VERSION  := $(shell git describe --tags)

LDFLAGS  := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -extldflags \"-static\""

.PHONY: $(/bin/bash egrep -o ^[a-zA-Z_-]+: $(MAKEFILE_LIST) | sed 's/://')

all: help

setup:	## setup dev tools
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint

test:	## run test
	go test -v -cover ./...

lint:	## run lint
	go vet ./...
	golint -set_exit_status *.go

build: deps ./cmd/wsbridge	## build binary
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) ./cmd/wsbridge

build-cross: deps	## build binary for each platforms
	for os in darwin linux windows; do \
	    for arch in amd64 386; do \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME) ./cmd/wsbridge; \
	    done; \
	done

install: test lint	## install binary
	go install $(LDFALGS)


deps:	## install dependencies
	dep ensure -v

up:	## update dependencies
	dep ensure -v -update

help:   ## show help
	@echo Usage: make [target]
	@echo ${\n}
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:		## clean bin, vendor
	go clean
	rm -rf bin/*
	rm -rf vendor/*
