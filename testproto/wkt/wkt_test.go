package wkt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/emptypb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/wrapperspb"
)

func TestWellKnownTypes(t *testing.T) {
	dur := durationpb.New(4*time.Hour + 2*time.Second)

	m := &MessageWithWKT{
		Duration:    dur,
		Empty:       &emptypb.Empty{},
		Timestamp:   timestamppb.Now(),
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

	// TODO prptoc-gen-prototxt
	// assert.Equal(t, golangMsg.String(), vtProtoMsg.String())
}
