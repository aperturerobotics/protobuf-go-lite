package test_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	"github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

var gogoMarshaler = jsonpb.Marshaler{
	OrigName:    true,
	EnumsAsInts: true,
}

func gogoMarshal(t *testing.T, msg proto.Message) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := gogoMarshaler.Marshal(&buf, msg); err != nil {
		t.Logf("gogo failed to marshal: %v", err)
	}
	return buf.Bytes()
}

var gogoUnmarshaler = jsonpb.Unmarshaler{}

func gogoUnmarshal(t *testing.T, msg proto.Message, data []byte) {
	t.Helper()
	buf := bytes.NewBuffer(data)
	if err := gogoUnmarshaler.Unmarshal(buf, msg); err != nil {
		t.Logf("gogo failed to unmarshal: %v", err)
	}
}

func generatedMarshal(t *testing.T, msg proto.Message, mask []string) []byte {
	t.Helper()
	m, ok := msg.(jsonplugin.Marshaler)
	if !ok {
		t.Fatalf("message %T does not implement the jsonplugin.Marshaler", msg)
	}
	s := jsonplugin.NewMarshalState(jsonplugin.DefaultMarshalerConfig).WithFieldMask(mask...)
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
	s := jsonplugin.NewUnmarshalState(data, jsonplugin.DefaultUnmarshalerConfig)
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

	gogoMarshaled := gogoMarshal(t, msg)

	generatedMarshaled := generatedMarshal(t, msg, mask)
	generatedFormatted := indent(t, generatedMarshaled)
	generatedDiff := cmp.Diff(expectedFormatted, generatedFormatted)

	if generatedDiff != "" {
		t.Errorf("expected : %s", string(expected))
		t.Errorf("gogo     : %s", string(gogoMarshaled))
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

	expectedMsgText := proto.MarshalTextString(msg)

	gogoUnmarshaled := reflect.New(reflect.ValueOf(msg).Elem().Type()).Interface().(proto.Message)
	gogoUnmarshal(t, gogoUnmarshaled, expected)
	gogoMsgText := proto.MarshalTextString(gogoUnmarshaled)

	generatedUnmarshaled := reflect.New(reflect.ValueOf(msg).Elem().Type()).Interface().(proto.Message)
	mask := generatedUnmarshal(t, generatedUnmarshaled, expected)
	generatedMsgText := proto.MarshalTextString(generatedUnmarshaled)
	generatedDiff := cmp.Diff(expectedMsgText, generatedMsgText)
	maskDiff := cmp.Diff(expectedMask, mask)

	if generatedDiff != "" {
		t.Errorf("expected : %s", string(expectedMsgText))
		t.Errorf("gogo     : %s", string(gogoMsgText))
		t.Errorf("generated: %s", string(generatedMsgText))
		if generatedDiff != "" {
			t.Errorf("  diff   : %s", generatedDiff)
		}
	}

	if maskDiff != "" {
		t.Errorf("mask diff: %s", maskDiff)
	}
}
