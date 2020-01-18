# Copyright 2019 Manlio Perillo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Exported variable definitions.
export GO111MODULE := on

# Imported variables: GO_PKG

# Variable definitions.
TESTFLAGS := -race -v
BENCHFLAGS := -v
COVERMODE := atomic # atomic is necessary if the -race flag is enabled

# Standard rules.
.POSIX:

.PHONY: build bench clean cover github install lint print test test-all trace vet

# Default rule.
build:
	go build

# Custom rules.
bench:
	go test ${BENCHFLAGS} -bench=. -benchmem ./...

clean:
	go mod tidy
	go clean
	go clean -i
	rm -f build/*

cover:
	go tool cover -html=build/coverage.out -o=build/coverage.html

github:
	git push --follow-tags -u github master

install:
	go install ./...

lint:
	golint ./...

print:
	goprint -font='"Inconsolata" 10pt/12pt' ${GO_PKG} > build/pkg.html
	prince -o build/pkg.pdf build/pkg.html

test:
	go test ${TESTFLAGS} -covermode=${COVERMODE} \
		-coverprofile=build/coverage.out \
		-trace=build/trace.out ${GO_PKG}

test-all:
	go test ${TESTFLAGS} ./...

trace:
	go tool trace build/trace.out

vet:
	go vet ./...
