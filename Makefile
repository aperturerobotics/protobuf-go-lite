# Copyright © 2021 The Things Industries B.V.
# SPDX-License-Identifier: Apache-2.0

.PHONY: default

default: build test

.PHONY: clean

clean:
	rm -f ./annotations/*.pb.go
	rm -f ./test/*/*.pb.go

.dev/protoc-gen-go-json/annotations.proto: annotations.proto
	mkdir -p $(shell dirname $@)
	cp $< $@

annotations/annotations.pb.go: .dev/protoc-gen-go-json/annotations.proto .dev/golangproto/bin/protoc .dev/golangproto/bin/protoc-gen-go
	PATH="$$PWD/.bin:$$PWD/.dev/golangproto/bin:$$PATH" protoc -I .dev --go_opt=module=github.com/TheThingsIndustries/protoc-gen-go-json --go_out=./ $<

BINARY_DEPS = annotations/annotations.pb.go $(wildcard cmd/protoc-gen-go-json/*.go) $(wildcard internal/gen/*.go)

VERSION ?= 0.0.0-dev

LDFLAGS = -X github.com/TheThingsIndustries/protoc-gen-go-json/internal/gen.Version=$(VERSION)

.bin/protoc-gen-go-json: $(BINARY_DEPS)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/protoc-gen-go-json

.bin/protoc-gen-go-json-linux-amd64: $(BINARY_DEPS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/protoc-gen-go-json

.bin/protoc-gen-go-json-linux-arm64: $(BINARY_DEPS)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/protoc-gen-go-json

.PHONY: build

build: .bin/protoc-gen-go-json .bin/protoc-gen-go-json-linux-amd64 .bin/protoc-gen-go-json-linux-arm64

.PHONY: watch

watch:
	ls annotations.proto cmd/protoc-gen-go-json/*.go internal/gen/*.go test/*.proto | entr make build test

OS :=
ifeq ($(shell uname),Linux)
	OS = linux
endif
ifeq ($(shell uname),Darwin)
	OS = osx
endif

.dev/golangproto/bin/protoc:
	mkdir -p .dev/golangproto/bin
	curl -sSL -o .dev/golangproto/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-$(OS)-x86_64.zip
	unzip -o .dev/golangproto/protoc.zip -d .dev/golangproto/

.dev/golangproto/bin/protoc-gen-go:
	go build -o $@ google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: testprotos

testprotos: build .dev/golangproto/bin/protoc .dev/golangproto/bin/protoc-gen-go
	PATH="$$PWD/.bin:$$PWD/.dev/golangproto/bin:$$PATH" protoc -I ./test -I . \
	  --go_opt=paths=source_relative --go_out=./test/golang \
	  --go-json_opt=paths=source_relative --go-json_opt=std=true --go-json_out=./test/golang \
	  ./test/*.proto

.PHONY: test

test: testprotos
	go test ./jsonplugin ./test/golang
