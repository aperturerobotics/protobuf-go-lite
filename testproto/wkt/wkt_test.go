package wkt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	"github.com/aperturerobotics/protobuf-go-lite/json"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/anypb"
	anypb_resolver "github.com/aperturerobotics/protobuf-go-lite/types/known/anypb/resolver"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/emptypb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/wrapperspb"
)

func TestWellKnownTypes(t *testing.T) {
	dur := durationpb.New(4*time.Hour + 2*time.Second)

	anyTypeResolver := anypb_resolver.NewFuncAnyTypeResolver(func(url string) (func() protobuf_go_lite.Message, error) {
		if url == "cool.apps/test-value" {
			return func() protobuf_go_lite.Message { return &durationpb.Duration{} }, nil
		}
		return nil, nil
	})
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
	mc := json.MarshalerConfig{
		AnyTypeResolver: anyTypeResolver,
	}
	jdata, err := mc.Marshal(m)
	if err != nil {
		require.NoError(t, err)
	}
	t.Log(string(jdata))

	// Ensure output is consistent
	var expected = `{"any":{"@type":"cool.apps/test-value","value":"14402s"},"duration":"14402s","empty":{},"timestamp":"2024-01-10T04:20:00Z","doubleValue":123456789.12345679,"floatValue":123456792,"int64Value":"123456789","uint64Value":"123456789","int32Value":"123456789","uint32Value":"123456789","boolValue":true,"stringValue":"String marshalling and unmarshalling test","bytesValue":"Qnl0ZXMgbWFyc2hhbGxpbmcgYW5kIHVubWFyc2hhbGxpbmcgdGVzdA=="}`
	require.Equal(t, expected, string(jdata))

	jparsed := &MessageWithWKT{}
	umc := json.UnmarshalerConfig{
		AnyTypeResolver: anyTypeResolver,
	}
	err = umc.Unmarshal(jdata, jparsed)
	require.NoError(t, err)
	require.True(t, jparsed.EqualVT(m))
}
