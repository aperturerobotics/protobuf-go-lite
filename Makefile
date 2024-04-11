SHELL:=bash
PROTOWRAP=hack/bin/protowrap
PROTOC_GEN_GO_LITE=hack/bin/protoc-gen-go-lite
PROTOC_GEN_VTPROTO_LITE=hack/bin/protoc-gen-go-lite-vtproto
GOIMPORTS=hack/bin/goimports
GOFUMPT=hack/bin/gofumpt
GOLANGCI_LINT=hack/bin/golangci-lint
GO_MOD_OUTDATED=hack/bin/go-mod-outdated
GOLIST=go list -f "{{ .Dir }}" -m

export GO111MODULE=on
undefine GOOS
undefine GOARCH

all:

vendor:
	go mod vendor

$(PROTOC_GEN_GO_LITE):
	cd ./hack; \
	go build -v \
		-o ./bin/protoc-gen-go-lite \
		github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite

$(PROTOC_GEN_VTPROTO_LITE):
	cd ./hack; \
	go build -v \
		-o ./bin/protoc-gen-go-lite-vtproto \
		github.com/aperturerobotics/vtprotobuf-lite/cmd/protoc-gen-go-lite-vtproto

$(GOIMPORTS):
	cd ./hack; \
	go build -v \
		-o ./bin/goimports \
		golang.org/x/tools/cmd/goimports

$(GOFUMPT):
	cd ./hack; \
	go build -v \
		-o ./bin/gofumpt \
		mvdan.cc/gofumpt

$(PROTOWRAP):
	cd ./hack; \
	go build -v \
		-o ./bin/protowrap \
		github.com/aperturerobotics/goprotowrap/cmd/protowrap

$(GOLANGCI_LINT):
	cd ./hack; \
	go build -v \
		-o ./bin/golangci-lint \
		github.com/golangci/golangci-lint/cmd/golangci-lint

$(GO_MOD_OUTDATED):
	cd ./hack; \
	go build -v \
		-o ./bin/go-mod-outdated \
		github.com/psampaz/go-mod-outdated

node_modules:
	yarn install

.PHONY: gen
gen: gen-wkt

PROTOBUF_ROOT=./lib/protobuf

.PHONY: gen-deps
gen-deps: vendor $(GOIMPORTS) $(PROTOWRAP) $(PROTOC) $(PROTOC_GEN_GO_LITE) $(PROTOC_GEN_VTPROTO_LITE)
	git submodule update --init ./lib/protobuf

.PHONY: gen-wkt
gen-wkt: gen-deps
	protoc \
		-I$(PROTOBUF_ROOT)/src \
		--plugin protoc-gen-go-lite-vtproto="$(PROTOC_GEN_VTPROTO_LITE)" \
		--go-lite-vtproto_out=. \
		--go-lite-vtproto_opt=module=google.golang.org/protobuf,wrap=true \
		--go-lite-vtproto_opt=module=github.com/aperturerobotics/protobuf-go-lite,wrap=true \
		--go-lite-vtproto_opt=Msrc/google/protobuf/descriptor.proto=github.com/aperturerobotics/vtprotobuf-lite/types/descriptorpb\;descriptorpb \
			$(PROTOBUF_ROOT)/src/google/protobuf/duration.proto \
			$(PROTOBUF_ROOT)/src/google/protobuf/descriptor.proto \
			$(PROTOBUF_ROOT)/src/google/protobuf/empty.proto \
			$(PROTOBUF_ROOT)/src/google/protobuf/timestamp.proto \
			$(PROTOBUF_ROOT)/src/google/protobuf/wrappers.proto \
			$(PROTOBUF_ROOT)/src/google/protobuf/struct.proto

.PHONY: outdated
outdated: $(GO_MOD_OUTDATED)
	go list -mod=mod -u -m -json all | $(GO_MOD_OUTDATED) -update -direct

.PHONY: list
list: $(GO_MOD_OUTDATED)
	go list -mod=mod -u -m -json all | $(GO_MOD_OUTDATED)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --timeout=10m

.PHONY: fix
fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix --timeout=10m

.PHONY: format
format: $(GOFUMPT) $(GOIMPORTS)
	$(GOIMPORTS) -w ./
	$(GOFUMPT) -w ./

.PHONY: test
test:
	go test -v ./...
