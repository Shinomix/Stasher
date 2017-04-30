NAME ?= stasher
PACKAGES = $$(go list ./... | grep -v 'vendor/')
SHELL := /bin/bash

init:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	govendor sync -v

fmt:
	@test -z "$$(gofmt -s -l . | grep -v vendor/)"

lint:
	@golint -set_exit_status $(PACKAGES)

test:
	@mkdir -p .coverage
	@for pkg in $(PACKAGES) ; do \
		go test $${pkg} -cover -coverprofile=.coverage/$${pkg//\//.}.part.out; \
		if [ $${?} -gt "0" ] ; then \
			exit 1; \
		fi \
	done
	@echo "mode: set" > .coverage/total.out
	@for f in `find . -name \*.part.out`; do tail -n +2 $$f >> .coverage/total.out; done
	@go tool cover -func=.coverage/total.out

build:
	@go build -o $(NAME)

.PHONY:
	init fmt lint test build
