package sizebaseline

import (
	"bytes"
	stdjson "encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/structpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/wrapperspb"
)

var (
	benchmarkSizeBaselineBool    bool
	benchmarkSizeBaselineJSON    []byte
	benchmarkSizeBaselineMessage *SizeBaseline
	benchmarkSizeBaselineSize    int
	benchmarkSizeBaselineString  string
	benchmarkSizeBaselineWire    []byte
)

func newSizeBaseline(tb testing.TB) *SizeBaseline {
	tb.Helper()

	structValue, err := structpb.NewStruct(map[string]any{
		"name":    "fixture",
		"enabled": true,
		"count":   3,
	})
	require.NoError(tb, err)

	valueValue, err := structpb.NewValue(map[string]any{
		"kind": "value",
	})
	require.NoError(tb, err)

	listValue, err := structpb.NewList([]any{"alpha", float64(2), true})
	require.NoError(tb, err)

	explicitInt32 := int32(11)
	explicitInt64 := int64(22)
	explicitUint32 := uint32(33)
	explicitUint64 := uint64(44)
	explicitSint32 := int32(-55)
	explicitSint64 := int64(-66)
	fixed32Value := uint32(77)
	fixed64Value := uint64(88)
	sfixed32Value := int32(-99)
	sfixed64Value := int64(-111)
	floatValue := float32(1.5)
	doubleValue := float64(2.5)
	boolValue := true
	stringValue := "baseline"
	requiredInt32 := int32(7)
	state := SizeBaseline_STATE_READY

	return &SizeBaseline{
		ExplicitInt32:  &explicitInt32,
		ImplicitInt32:  12,
		ExplicitInt64:  &explicitInt64,
		ExplicitUint32: &explicitUint32,
		ExplicitUint64: &explicitUint64,
		ExplicitSint32: &explicitSint32,
		ExplicitSint64: &explicitSint64,
		Fixed32Value:   &fixed32Value,
		Fixed64Value:   &fixed64Value,
		Sfixed32Value:  &sfixed32Value,
		Sfixed64Value:  &sfixed64Value,
		FloatValue:     &floatValue,
		DoubleValue:    &doubleValue,
		BoolValue:      &boolValue,
		StringValue:    &stringValue,
		BytesValue:     []byte("bytes"),
		RequiredInt32:  &requiredInt32,
		PackedInt32:    []int32{1, 2, 3},
		ExpandedInt32:  []int32{4, 5, 6},
		NestedValues: []*SizeBaseline_Nested{
			{Name: ptrString("first"), Count: 1, Labels: []string{"a", "b"}},
			{Name: ptrString("second"), Count: 2, Labels: []string{"c"}},
		},
		NestedByName: map[string]*SizeBaseline_Nested{
			"primary": {Name: ptrString("primary"), Count: 10, Labels: []string{"map"}},
		},
		NestedById: map[uint32]*SizeBaseline_Nested{
			1: {Name: ptrString("one"), Count: 1},
			2: {Name: ptrString("two"), Count: 2},
		},
		State:         &state,
		Nested:        &SizeBaseline_Nested{Name: ptrString("nested"), Count: 99, Labels: []string{"x"}},
		Timestamp:     timestamppb.New(time.Date(2026, time.May, 30, 12, 34, 56, 0, time.UTC)),
		Duration:      durationpb.New(5*time.Second + 250*time.Millisecond),
		StringWrapper: wrapperspb.String("wrapped"),
		BytesWrapper:  wrapperspb.Bytes([]byte("wrapped-bytes")),
		StructValue:   structValue,
		ValueValue:    valueValue,
		ListValue:     listValue,
		Selection:     &SizeBaseline_SelectedName{SelectedName: "chosen"},
	}
}

func ptrString(v string) *string {
	return &v
}

func TestSizeBaselineGeneratedBehavior(t *testing.T) {
	msg := newSizeBaseline(t)

	wire, err := msg.MarshalVT()
	require.NoError(t, err)
	require.NotEmpty(t, wire)
	require.Equal(t, len(wire), msg.SizeVT())

	var out SizeBaseline
	require.NoError(t, out.UnmarshalVT(wire))
	require.True(t, out.EqualVT(msg))
	require.True(t, msg.CloneVT().EqualVT(msg))

	jsonMsg := msg.CloneVT()
	jsonMsg.StructValue = nil
	jsonMsg.ValueValue = nil
	jsonMsg.ListValue = nil

	jsonData, err := stdjson.Marshal(jsonMsg)
	require.NoError(t, err)
	require.Contains(t, string(jsonData), `"selectedName":"chosen"`)

	var fromJSON SizeBaseline
	require.NoError(t, stdjson.Unmarshal(jsonData, &fromJSON))
	require.Equal(t, jsonMsg.GetRequiredInt32(), fromJSON.GetRequiredInt32())
	require.Equal(t, jsonMsg.GetStringValue(), fromJSON.GetStringValue())
	require.Equal(t, jsonMsg.GetSelectedName(), fromJSON.GetSelectedName())
	require.Equal(t, jsonMsg.GetTimestamp().GetSeconds(), fromJSON.GetTimestamp().GetSeconds())

	textData := msg.String()
	for _, expected := range []string{
		`string_value: "baseline"`,
		"packed_int32: [1, 2, 3]",
		`selected_name: "chosen"`,
		"timestamp: seconds:",
	} {
		require.Truef(t, strings.Contains(textData, expected), "text output missing %q: %s", expected, textData)
	}
}

func TestSizeBaselineUnknownFieldsAndRequired(t *testing.T) {
	_, err := (&SizeBaseline{}).MarshalVT()
	require.ErrorContains(t, err, "required field required_int32 not set")

	wire, err := newSizeBaseline(t).MarshalVT()
	require.NoError(t, err)

	unknownField := []byte{0x98, 0x06, 0x7b}
	wireWithUnknown := append(append([]byte{}, wire...), unknownField...)

	var out SizeBaseline
	require.NoError(t, out.UnmarshalVT(wireWithUnknown))
	require.True(t, out.CloneVT().EqualVT(&out))

	remarshaled, err := out.MarshalVT()
	require.NoError(t, err)
	require.True(t, bytes.Contains(remarshaled, unknownField), "unknown field was not preserved")
	require.Equal(t, len(remarshaled), out.SizeVT())
}

func BenchmarkSizeBaselineMarshalVT(b *testing.B) {
	msg := newSizeBaseline(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wire, err := msg.MarshalVT()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkSizeBaselineWire = wire
	}
}

func BenchmarkSizeBaselineUnmarshalVT(b *testing.B) {
	wire, err := newSizeBaseline(b).MarshalVT()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg SizeBaseline
		if err := msg.UnmarshalVT(wire); err != nil {
			b.Fatal(err)
		}
		benchmarkSizeBaselineMessage = &msg
	}
}

func BenchmarkSizeBaselineMarshalJSON(b *testing.B) {
	msg := newSizeBaseline(b).CloneVT()
	msg.StructValue = nil
	msg.ValueValue = nil
	msg.ListValue = nil

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := stdjson.Marshal(msg)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkSizeBaselineJSON = data
	}
}

func BenchmarkSizeBaselineUnmarshalJSON(b *testing.B) {
	msg := newSizeBaseline(b).CloneVT()
	msg.StructValue = nil
	msg.ValueValue = nil
	msg.ListValue = nil
	data, err := stdjson.Marshal(msg)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out SizeBaseline
		if err := stdjson.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
		benchmarkSizeBaselineMessage = &out
	}
}

func BenchmarkSizeBaselineMarshalProtoText(b *testing.B) {
	msg := newSizeBaseline(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSizeBaselineString = msg.String()
	}
}

func BenchmarkSizeBaselineCloneVT(b *testing.B) {
	msg := newSizeBaseline(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSizeBaselineMessage = msg.CloneVT()
	}
}

func BenchmarkSizeBaselineEqualVT(b *testing.B) {
	msg := newSizeBaseline(b)
	clone := msg.CloneVT()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSizeBaselineBool = msg.EqualVT(clone)
	}
}

func BenchmarkSizeBaselineSizeVT(b *testing.B) {
	msg := newSizeBaseline(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSizeBaselineSize = msg.SizeVT()
	}
}
