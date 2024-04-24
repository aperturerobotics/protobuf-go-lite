package anypb

import (
	"errors"
	fmt "fmt"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	"github.com/aperturerobotics/protobuf-go-lite/json"
	anypb_resolver "github.com/aperturerobotics/protobuf-go-lite/types/known/anypb/resolver"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/structpb"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
)

// MessageTypeResolver is an interface for looking up messages.
//
// A compliant implementation must deterministically return the same type
// if no error is encountered.
//
// The [Types] type implements this interface.
type MessageTypeResolver = anypb_resolver.AnyTypeResolver

// ErrNotFound is returned if the message type was not found.
var ErrNotFound = errors.New("message type not found in Any")

// New marshals src into a new Any instance.
func New(src protobuf_go_lite.Message, typeURL string) (*Any, error) {
	dst := new(Any)
	if err := dst.MarshalFrom(src, typeURL); err != nil {
		return nil, err
	}
	return dst, nil
}

// MarshalFrom marshals src into dst as the underlying message
// using the provided marshal options.
//
// If no options are specified, call dst.MarshalFrom instead.
func MarshalFrom(dst *Any, src protobuf_go_lite.Message, typeURL string) error {
	if src == nil {
		dst.Reset()
		return nil
	}
	b, err := src.MarshalVT()
	if err != nil {
		return err
	}
	dst.TypeUrl = typeURL
	dst.Value = b
	return nil
}

// UnmarshalTo unmarshals the underlying message from src into dst
// using the provided unmarshal options.
// It reports an error if dst is not of the right message type.
//
// If no options are specified, call src.UnmarshalTo instead.
func UnmarshalTo(src *Any, dst protobuf_go_lite.Message, typeURL string) error {
	if src == nil {
		dst.Reset()
		return nil
	}
	if !src.MessageIs(typeURL) {
		got := typeURL
		want := src.GetTypeUrl()
		return fmt.Errorf("mismatched message type: got %q, want %q", got, want)
	}
	return dst.UnmarshalVT(src.GetValue())
}

// UnmarshalNew unmarshals the underlying message from src into dst,
// which is newly created message using a type resolved from the type URL.
// The message type is resolved according to opt.Resolver,
// which should implement protoregistry.MessageTypeResolver.
// It reports an error if the underlying message type could not be resolved.
//
// If no options are specified, call src.UnmarshalNew instead.
func UnmarshalNew(src *Any, typeURL string, resolver MessageTypeResolver) (dst protobuf_go_lite.Message, err error) {
	if src.GetTypeUrl() == "" {
		return nil, errors.New("invalid empty type URL")
	}
	if resolver == nil {
		return nil, errors.New("message type resolver cannot be empty")
	}
	mt, err := resolver.FindMessageByURL(src.GetTypeUrl())
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not resolve %q: %w", src.GetTypeUrl(), err)
	}
	dst = mt()
	if dst == nil {
		return nil, ErrNotFound
	}
	return dst, nil
}

// MessageIs reports whether the underlying message is of the same type as m.
func (x *Any) MessageIs(typeURL string) bool {
	return x.GetTypeUrl() == typeURL
}

// MarshalFrom marshals m into x as the underlying message.
func (x *Any) MarshalFrom(m protobuf_go_lite.Message, typeURL string) error {
	return MarshalFrom(x, m, typeURL)
}

// UnmarshalTo unmarshals the contents of the underlying message of x into m.
// It resets m before performing the unmarshal operation.
// It reports an error if m is not of the right message type.
func (x *Any) UnmarshalTo(m protobuf_go_lite.Message, typeURL string) error {
	return UnmarshalTo(x, m, typeURL)
}

// UnmarshalNew unmarshals the contents of the underlying message of x into
// a newly allocated message of the specified type.
// It reports an error if the underlying message type could not be resolved.
func (x *Any) UnmarshalNew(typeURL string, resolver MessageTypeResolver) (protobuf_go_lite.Message, error) {
	return UnmarshalNew(x, typeURL, resolver)
}

// MarshalJSON marshals the Any to JSON.
func (x *Any) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Any from JSON.
func (x *Any) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals an Any WKT.
// UnmarshalProtoJSON unmarshals an Any WKT.
func (x *Any) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}

	// Read the raw object and create a sub-unmarshaler for it.
	data := s.SkipAndReturnBytes()
	if s.Err() != nil {
		return
	}
	sub := s.Sub(data)

	// Read the first field in the object. This should be @type.
	if key := sub.ReadObjectField(); key != "@type" {
		s.SetErrorf("first field in Any is not @type, but %q", key)
		return
	}
	typeURL := sub.ReadString()
	if err := sub.Err(); err != nil {
		return
	}

	// Find the message type by the type URL.
	t, err := s.AnyTypeResolver().FindMessageByURL(typeURL)
	if err != nil {
		s.SetError(err)
		return
	}

	// Allocate a new message of that type.
	msg := t()
	if msg == nil {
		s.SetError(ErrNotFound)
		return
	}

	var unmarshaler json.Unmarshaler
	switch msgt := msg.(type) {
	default:
		s.SetError(ErrNotFound)
		return
	case *durationpb.Duration,
		*structpb.Struct,
		*structpb.Value,
		*structpb.ListValue,
		*timestamppb.Timestamp:
		if field := sub.ReadObjectField(); field != "value" {
			s.SetErrorf("unexpected %q field in Any", field)
			return
		}
		unmarshaler = msgt.(json.Unmarshaler)
	case json.Unmarshaler:
		// Create another sub-unmarshaler for the raw data and unmarshal the message.
		sub = s.Sub(data)
	}

	unmarshaler.UnmarshalProtoJSON(sub)
	if err := sub.Err(); err != nil {
		return
	}

	if field := sub.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in Any", field)
		return
	}

	// Wrap the unmarshaled message in an Any and return that.
	n, err := New(msg, typeURL)
	if err != nil {
		sub.SetError(err)
		return
	} else {
		*x = *n
	}
}

// MarshalProtoJSON marshals an Any WKT.
func (x *Any) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}

	s.WriteObjectStart()
	s.WriteObjectField("@type")
	s.WriteString(x.GetTypeUrl())

	mt, err := s.AnyTypeResolver().FindMessageByURL(x.GetTypeUrl())
	if err != nil {
		s.SetError(err)
		return
	}

	msg := mt()
	if err := x.UnmarshalTo(msg, x.GetTypeUrl()); err != nil {
		s.SetError(err)
		return
	}

	s.WriteMore()
	s.WriteObjectField("value")

	switch m := msg.(type) {
	case json.Marshaler:
		m.MarshalProtoJSON(s)
	default:
		s.SetError(errors.New("message in Any does not implement json.Marshaler"))
	}

	s.WriteObjectEnd()
}
