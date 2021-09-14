# protoc-gen-go-json

> Protoc plugin for generating JSON marshalers and unmarshalers in Go

## Background

The API of The Things Stack V3 is not compliant with the latest `protojson` package because we used github.com/gogo/protobuf with custom types and custom JSON marshalers. We do want to upgrade to the V2 protobuf library, but also not break the API of The Things Stack. Therefore we'll now generate our own JSON marshalers and unmarshalers.

## Usage in Proto Code

```proto
syntax = "proto3";

import "github.com/TheThingsIndustries/protoc-gen-go-json/annotations.proto";

package thethings.json.example;

option go_package = "github.com/TheThingsIndustries/protoc-gen-go-json/example";

option (thethings.json.file) = {
  marshaler_all: true,   // Generate marshalers for everything in the file.
  unmarshaler_all: true, // Generate unmarshalers for everything in the file.
};

// This is an enum with custom JSON marshaling.
enum MyCustomEnum {
  option (thethings.json.enum) = {
    marshal_as_string: true,  // The marshaler will render values as strings.
    prefix: "CUSTOM"          // The unmarshaler will accept both UNKNOWN and CUSTOM_UNKNOWN, etc.
  };

  CUSTOM_UNKNOWN = 0;
  CUSTOM_V1_0 = 1 [
    (thethings.json.enum_value) = {
      value: "1.0",            // The marshaler will render this value as "1.0".
      aliases: ["1.0.0"]       // The unmarshaler will also accept "1.0.0".
    }
  ];
  CUSTOM_V1_0_1 = 2 [
    (thethings.json.enum_value) = {
      value: "1.0.1"           // The marshaler will render this value as "1.0.1".
    }
  ];
}

message CustomEnumWrapper {
  option (thethings.json.message) = {
    wrapper: true                // This message acts as a wrapper WKT.
  };
  CustomEnum value = 1;
}

message MessageWithCustomRenderedField {
  bytes hex_field = 1 [
    (thethings.json.field) = {
      marshaler_func: "github.com/TheThingsIndustries/protoc-gen-go-json/test/types.MarshalHEX",
      unmarshaler_func: "github.com/TheThingsIndustries/protoc-gen-go-json/test/types.UnmarshalHEX"
    }
  ];
}
```

## Generating Go Code

```bash
$ protoc -I ./path/to -I . \
  --go_opt=paths=source_relative --go_out=./path/to \
  --go-json_opt=paths=source_relative --go-json_out=./path/to \
  ./path/to/*.proto
```

## Usage in Go Code

```go
data, err := jsonplugin.MarshalerConfig{
	// config ...
}.Marshal(msg)
```

```go
err := jsonplugin.UnmarshalerConfig{
	// config...
}.Unmarshal(data, msg)
```

## Contributing

We do not accept external issues with feature requests. This plugin only supports what we actually use ourselves at The Things Industries.

We do not accept external pull requests with new features, but everyone is free to fork it and add features in their own fork.

We do accept issues with bug reports and pull requests with bug fixes.
