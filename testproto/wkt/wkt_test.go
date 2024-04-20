package wkt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aperturerobotics/protobuf-go-lite/types/known/anypb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/emptypb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/wrapperspb"
)

func TestWellKnownTypes(t *testing.T) {
	dur := durationpb.New(4*time.Hour + 2*time.Second)

	anyVal, err := anypb.New(dur, "cool.apps/test-value")
	require.NoError(t, err)

	ts := timestamppb.New(time.Date(2024, time.January, 10, 4, 20, 00, 00, time.UTC))
	m := &MessageWithWKT{
		Any:         anyVal,
		Duration:    dur,
		Empty:       &emptypb.Empty{},
		Timestamp:   ts,
		DoubleValue: wrapperspb.Double(123456789.123456789),
		FloatValue:  wrapperspb.Float(123456789.123456789),
		Int64Value:  wrapperspb.Int64(123456789),
		Uint64Value: wrapperspb.UInt64(123456789),
		Int32Value:  wrapperspb.Int32(123456789),
		Uint32Value: wrapperspb.UInt32(123456789),
		BoolValue:   wrapperspb.Bool(true),
		StringValue: wrapperspb.String("String marshalling and unmarshalling test"),
		BytesValue:  wrapperspb.Bytes([]byte("Bytes marshalling and unmarshalling test")),
	}

	vtProtoBytes, err := m.MarshalVT()
	require.NoError(t, err)

	require.NotEmpty(t, vtProtoBytes)

	var (
		// golangMsg  = &MessageWithWKT{}
		vtProtoMsg = &MessageWithWKT{}
	)

	require.NoError(t, vtProtoMsg.UnmarshalVT(vtProtoBytes))

	// TODO prototxt
	// assert.Equal(t, golangMsg.String(), vtProtoMsg.String())

	// TODO protoc json
	jdata, err := m.MarshalJSON()
	if err != nil {
		require.NoError(t, err)
	}
	t.Log(string(jdata))

	// Ensure output is consistent
	var expected = `{"any":{"type_url":"cool.apps/test-value","value":"CMJw"},"duration":{"seconds":"14402"},"empty":{},"timestamp":{"seconds":"1704860400"},"double_value":{"value":123456789.12345679},"float_value":{"value":123456790},"int64_value":{"value":"123456789"},"uint64_value":{"value":"123456789"},"int32_value":{"value":123456789},"uint32_value":{"value":123456789},"bool_value":{"value":true},"string_value":{"value":"String marshalling and unmarshalling test"},"bytes_value":{"value":"Qnl0ZXMgbWFyc2hhbGxpbmcgYW5kIHVubWFyc2hhbGxpbmcgdGVzdA=="}}`
	require.Equal(t, expected, string(jdata))
}
