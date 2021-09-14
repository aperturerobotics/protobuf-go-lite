package test_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	proto "google.golang.org/protobuf/proto"
)

var pluginMarshaler = jsonplugin.MarshalerConfig{
	EnumsAsInts: true,
}

func generatedMarshal(t *testing.T, msg proto.Message, mask []string) []byte {
	t.Helper()
	m, ok := msg.(jsonplugin.Marshaler)
	if !ok {
		t.Fatalf("message %T does not implement the jsonplugin.Marshaler", msg)
	}
	s := jsonplugin.NewMarshalState(pluginMarshaler).WithFieldMask(mask...)
	m.MarshalProtoJSON(s)
	b, err := s.Bytes()
	if err != nil {
		t.Fatalf("generated failed to marshal: %v", err)
	}
	return b
}

func generatedUnmarshal(t *testing.T, msg proto.Message, data []byte) []string {
	t.Helper()
	unmarshaler, ok := msg.(jsonplugin.Unmarshaler)
	if !ok {
		t.Fatalf("message %T does not implement the jsonplugin.Unmarshaler", msg)
	}
	s := jsonplugin.NewUnmarshalState(data, jsonplugin.UnmarshalerConfig{})
	unmarshaler.UnmarshalProtoJSON(s)
	if err := s.Err(); err != nil {
		t.Fatalf("generated failed to unmarshal: %v", err)
	}
	paths := s.FieldMask().GetPaths()
	if len(paths) == 0 {
		return nil
	}
	return paths
}

func indent(t *testing.T, data []byte) string {
	t.Helper()
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		t.Fatalf("failed to indent %s: %v", string(data), err)
	}
	return buf.String()
}

func expectMarshalEqual(t *testing.T, msg proto.Message, mask []string, expected []byte) {
	t.Helper()

	expectedFormatted := indent(t, expected)

	generatedMarshaled := generatedMarshal(t, msg, mask)
	generatedFormatted := indent(t, generatedMarshaled)
	generatedDiff := cmp.Diff(expectedFormatted, generatedFormatted)

	if generatedDiff != "" {
		t.Errorf("expected : %s", string(expected))
		t.Errorf("generated: %s", string(generatedMarshaled))
		if generatedDiff != "" {
			t.Errorf("  diff   : %s", generatedDiff)
		}
	}
}

func expectUnmarshalEqual(t *testing.T, msg proto.Message, expected []byte, expectedMask []string) {
	t.Helper()
	if msg == nil {
		return
	}

	expectedMsgText := prototext.Format(msg)

	generatedUnmarshaled := reflect.New(reflect.ValueOf(msg).Elem().Type()).Interface().(proto.Message)
	mask := generatedUnmarshal(t, generatedUnmarshaled, expected)
	generatedMsgText := prototext.Format(generatedUnmarshaled)
	generatedDiff := cmp.Diff(expectedMsgText, generatedMsgText)
	maskDiff := cmp.Diff(expectedMask, mask)

	if generatedDiff != "" {
		t.Errorf("expected : %s", string(expectedMsgText))
		t.Errorf("generated: %s", string(generatedMsgText))
		if generatedDiff != "" {
			t.Errorf("  diff   : %s", generatedDiff)
		}
	}

	if maskDiff != "" {
		t.Errorf("mask diff: %s", maskDiff)
	}
}
