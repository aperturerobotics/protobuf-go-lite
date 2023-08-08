# Copyright © 2021 The Things Industries B.V.
# SPDX-License-Identifier: Apache-2.0

.PHONY: default

default: build test

.PHONY: clean

clean:
	rm -f ./annotations/*.pb.go
	rm -f ./test/*/*.pb.go

.bin/protoc-gen-go: go.mod
	GOBIN=$(PWD)/.bin go install google.golang.org/protobuf/cmd/protoc-gen-go

annotations/annotations.pb.go: api/thethings/json/annotations.proto .bin/protoc-gen-go
	buf generate api

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
	ls api/thethings/json/annotations.proto cmd/protoc-gen-go-json/*.go internal/gen/*.go test/*.proto | entr make build test

.PHONY: testprotos

testprotos: build .bin/protoc-gen-go
	buf generate --template buf.gen.test.yaml test

.PHONY: test

test: testprotos
	go test ./jsonplugin ./test/golang
