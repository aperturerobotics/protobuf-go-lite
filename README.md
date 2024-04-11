# protobuf-go-lite

[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/protobuf-go-lite.svg)](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite)
[![Build Status](https://travis-ci.org/protocolbuffers/protobuf-go.svg?branch=master)](https://travis-ci.org/protocolbuffers/protobuf-go)

## Introduction

This is a version of [protobuf-go] which does not use reflection.

[protobuf-go]: https://github.com/protocolbuffers/protobuf-go

Use it with [vtprotobuf-lite] to compile static code for marshaling,
unmarshaling, and sizing protobuf messages.

[vtprotobuf-lite]: https://github.com/aperturerobotics/vtprotobuf-lite

Use it with [protoc-gen-go-lite-json] to compile static code for json marshaling
and unmarshaling to avoid importing encoding/json and the reflect package.

[protoc-gen-go-lite-json]: https://github.com/aperturerobotics/protoc-gen-go-lite-json

This is both more efficient at runtime, generates smaller code binaries without
the reflection information, and features better support for [tinygo] which lacks
support for some reflection features.

[tinygo]: https://github.com/tinygo-org/tinygo

**This is a fork of the [upstream project] modified to not use reflection.**

[upstream project]: https://github.com/protocolbuffers/protobuf-go

## Protobuf

[protocol buffers](https://protobuf.dev) are a cross-platform cross-language
message serialization format. Protobuf is a language for specifying the schema
for structured data. This schema is compiled into language specific bindings.
This project provides both a tool to generate Go code for the protocol buffer
language, and also the runtime implementation to handle serialization of
messages in Go.

See the [protocol buffer developer guide](https://protobuf.dev/overview) for
more information about protocol buffers themselves.

## Example

See the [protobuf-project](https://github.com/aperturerobotics/protobuf-project)
template for an example of how to use this package and vtprotobuf together with
protowrap to generate protobufs for your project.

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

## License

BSD-3
