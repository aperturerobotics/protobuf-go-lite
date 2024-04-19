# protobuf-go-lite

[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/protobuf-go-lite.svg)](https://pkg.go.dev/github.com/aperturerobotics/protobuf-go-lite)
[![Build Status](https://travis-ci.org/protocolbuffers/protobuf-go.svg?branch=master)](https://travis-ci.org/protocolbuffers/protobuf-go)

## Introduction

**protobuf-go-lite** is a stripped-down version of the [protobuf-go] code
generator modified to work without reflection and merged with [vtprotobuf] to
provide modular features with static code generation for marshal/unmarshal,
size, clone, and equal.

[protobuf-go]: https://github.com/protocolbuffers/protobuf-go
[vtprotobuf]: https://github.com/planetscale/vtprotobuf

Static code generation without reflection is more efficient at runtime and
results in smaller code binaries. It also provides better support for [tinygo]
which has limited reflection support.

[tinygo]: https://github.com/tinygo-org/tinygo

Lightweight Protobuf 3 RPCs are implemented in [StaRPC] for Go and TypeScript.

[StaRPC]: https://github.com/aperturerobotics/starpc

[protoc-gen-doc] is recommended for generating documentation.

[protoc-gen-doc]: https://github.com/pseudomuto/protoc-gen-doc

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
        protoc \
          --plugin protoc-gen-go-lite="${GOBIN}/protoc-gen-go-lite"
          --go-lite_out=.  \
          --go-lite_opt=features=marshal+unmarshal+size+equal+clone \
        proto/$${name}.proto; \
    done
    ```

`protobuf-go-lite` replaces `protoc-gen-go` and `protoc-gen-go-vtprotobuf` and should not be used with those generators.

## Available features

The following additional features from vtprotobuf can be enabled:

- `size`: generates a `func (p *YourProto) SizeVT() int` helper that behaves identically to calling `proto.Size(p)` on the message, except the size calculation is fully unrolled and does not use reflection. This helper function can be used directly, and it'll also be used by the `marshal` codegen to ensure the destination buffer is properly sized before ProtoBuf objects are marshalled to it.

- `equal`: generates the following helper methods

    - `func (this *YourProto) EqualVT(that *YourProto) bool`: this function behaves almost identically to calling `proto.Equal(this, that)` on messages, except the equality calculation is fully unrolled and does not use reflection. This helper function can be used directly.

    - `func (this *YourProto) EqualMessageVT(thatMsg any) bool`: this function behaves like the above `this.EqualVT(that)`, but allows comparing against arbitrary proto messages. If `thatMsg` is not of type `*YourProto`, false is returned. The uniform signature provided by this method allows accessing this method via type assertions even if the message type is not known at compile time. This allows implementing a generic `func EqualVT(proto.Message, proto.Message) bool` without reflection.

- `marshal`: generates the following helper methods

    - `func (p *YourProto) MarshalVT() ([]byte, error)`: this function behaves identically to calling `proto.Marshal(p)`, except the actual marshalling has been fully unrolled and does not use reflection or allocate memory. This function simply allocates a properly sized buffer by calling `SizeVT` on the message and then uses `MarshalToSizedBufferVT` to marshal to it.

    - `func (p *YourProto) MarshalToVT(data []byte) (int, error)`: this function can be used to marshal a message to an existing buffer. The buffer must be large enough to hold the marshalled message, otherwise this function will panic. It returns the number of bytes marshalled. This function is useful e.g. when using memory pooling to re-use serialization buffers.

    - `func (p *YourProto) MarshalToSizedBufferVT(data []byte) (int, error)`: this function behaves like `MarshalTo` but expects that the input buffer has the exact size required to hold the message, otherwise it will panic.

- `marshal_strict`: generates the following helper methods

    - `func (p *YourProto) MarshalVTStrict() ([]byte, error)`: this function behaves like `MarshalVT`, except fields are marshalled in a strict order by field's numbers they were declared in .proto file.

    - `func (p *YourProto) MarshalToVTStrict(data []byte) (int, error)`: this function behaves like `MarshalToVT`, except fields are marshalled in a strict order by field's numbers they were declared in .proto file.

    - `func (p *YourProto) MarshalToSizedBufferVTStrict(data []byte) (int, error)`: this function behaves like `MarshalToSizedBufferVT`, except fields are marshalled in a strict order by field's numbers they were declared in .proto file.


- `unmarshal`: generates a `func (p *YourProto) UnmarshalVT(data []byte)` that behaves similarly to calling `proto.Unmarshal(data, p)` on the message, except the unmarshalling is performed by unrolled codegen without using reflection and allocating as little memory as possible. If the receiver `p` is **not** fully zeroed-out, the unmarshal call will actually behave like `proto.Merge(data, p)`. This is because the `proto.Unmarshal` in the ProtoBuf API is implemented by resetting the destination message and then calling `proto.Merge` on it. To ensure proper `Unmarshal` semantics, ensure you've called `proto.Reset` on your message before calling `UnmarshalVT`, or that your message has been newly allocated.

- `unmarshal_unsafe` generates a `func (p *YourProto) UnmarshalVTUnsafe(data []byte)` that behaves like `UnmarshalVT`, except it unsafely casts slices of data to `bytes` and `string` fields instead of copying them to newly allocated arrays, so that it performs less allocations. **Data received from the wire has to be left untouched for the lifetime of the message.** Otherwise, the message's `bytes` and `string` fields can be corrupted.

- `clone`: generates the following helper methods

    - `func (p *YourProto) CloneVT() *YourProto`: this function behaves similarly to calling `proto.Clone(p)` on the message, except the cloning is performed by unrolled codegen without using reflection. If the receiver `p` is `nil` a typed `nil` is returned.

    - `func (p *YourProto) CloneMessageVT() any`: this function behaves like the above `p.CloneVT()`, but provides a uniform signature in order to be accessible via type assertions even if the type is not known at compile time. This allows implementing a generic `func CloneMessageVT() any` without reflection. If the receiver `p` is `nil`, a typed `nil` pointer of the message type will be returned inside a `any` interface.

- `json`: generates the following helper methods

    - `func (p *YourProto) UnmarshalJSON(data []byte) error` behaves similarly to calling `protojson.Unmarshal(data, p)` on the message, except the unmarshalling is performed by unrolled codegen without using reflection and allocating as little memory as possible (with valyala/fastjson). If the receiver `p` is **not** fully zeroed-out, the unmarshal call will actually behave like `proto.Merge(data, p)`. To ensure proper `Unmarshal` semantics, ensure you've called `proto.Reset` on your message before calling `UnmarshalJSON`, or that your message has been newly allocated.

    - `func (p *YourProto) UnmarshalJSONValue(val *fastjson.Value) error` unmarshals a `*fastjson.Value`.

    - `func (p *YourProto) MarshalJSON() ([]byte, error)` behaves similarly to calling `protojson.Marshal(p)` on the message, except the marshalling is performed by unrolled codegen without using reflection and allocating as little memory as possible (with Jeffail/gabs).

## License

BSD-3
