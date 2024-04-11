# Go support for Protocol Buffers without Reflect

[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/protobuf-go-lite.svg)](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite)
[![Build Status](https://travis-ci.org/protocolbuffers/protobuf-go.svg?branch=master)](https://travis-ci.org/protocolbuffers/protobuf-go)

This project hosts the Go implementation for
[protocol buffers](https://protobuf.dev), which is a
language-neutral, platform-neutral, extensible mechanism for serializing
structured data. The protocol buffer language is a language for specifying the
schema for structured data. This schema is compiled into language specific
bindings. This project provides both a tool to generate Go code for the
protocol buffer language, and also the runtime implementation to handle
serialization of messages in Go. See the
[protocol buffer developer guide](https://protobuf.dev/overview)
for more information about protocol buffers themselves.

**This is a fork of the [upstream project] modified to not use reflection.**

[upstream project]: https://github.com/protocolbuffers/protobuf-go

This project is a Code generator: The
[`protoc-gen-go-lite`](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite)
tool is a compiler plugin to `protoc`, the protocol buffer compiler. It augments
the `protoc` compiler so that it knows how to [generate Go specific code for a
given `.proto` file](https://protobuf.dev/reference/go/go-generated).

It is recommended to also use [vtprotobuf-lite] to compile static code for
marshaling, unmarshaling, and sizing protobuf messages.

[vtprotobuf-lite]: https://github.com/aperturerobotics/vtprotobuf-lite

See the [protobuf-project](https://github.com/aperturerobotics/protobuf-project)
template for an example of how to use this package and vtprotobuf together with
protowrap to generate protobufs for your project.

See the
[developer guide for protocol buffers in Go](https://protobuf.dev/getting-started/gotutorial)
for a general guide for how to get started using protobufs in Go.

This package is available at **github.com/aperturerobotics/protobuf-go-lite**.

## Package index

Summary of the packages provided by this module:

*   [`compiler/protogen`](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite/compiler/protogen):
    Package `protogen` provides support for writing protoc plugins.
*   [`cmd/protoc-gen-go-lite`](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite):
    The `protoc-gen-go-lite` binary is a protoc plugin to generate a Go protocol
    buffer package.

## Usage

1. Install `protoc-gen-go-lite`:

    ```
    go install github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite@latest
    ```

2. Ensure your project is already using the ProtoBuf v2 API (i.e. `google.golang.org/protobuf`). The `vtprotobuf` compiler is not compatible with APIv1 generated code.

3. Update your `protoc` generator to use the new plug-in.

    ```
    for name in $(PROTO_SRC_NAMES); do \
        $(VTROOT)/bin/protoc \
        --go-lite_out=. --plugin protoc-gen-go-lite="${GOBIN}/protoc-gen-go-lite" \
        proto/$${name}.proto; \
    done
    ```
