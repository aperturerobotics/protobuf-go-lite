// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package jsonplugin

import jsoniter "github.com/json-iterator/go"

var jsoniterConfig = jsoniter.Config{
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func valueTypeString(v jsoniter.ValueType) string {
	switch v {
	case jsoniter.StringValue:
		return "String"
	case jsoniter.NumberValue:
		return "Number"
	case jsoniter.NilValue:
		return "Null"
	case jsoniter.BoolValue:
		return "Bool"
	case jsoniter.ArrayValue:
		return "Array"
	case jsoniter.ObjectValue:
		return "Object"
	default:
		return "unknown"
	}
}
