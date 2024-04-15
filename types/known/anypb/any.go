package anypb

import (
	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	"github.com/pkg/errors"
)

// MessageTypeResolver is an interface for looking up messages.
//
// A compliant implementation must deterministically return the same type
// if no error is encountered.
//
// The [Types] type implements this interface.
type MessageTypeResolver interface {
	// FindMessageByURL looks up a message by a URL identifier.
	// See documentation on google.protobuf.Any.type_url for the URL format.
	//
	// Returns the constructor for the message.
	// This returns (nil, ErrNotFound) if not found.
	FindMessageByURL(url string) (func() protobuf_go_lite.Message, error)
}

// ErrNotFound is returned if the message type was not found.
var ErrNotFound = errors.New("proto type not found")

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
		return errors.Errorf("mismatched message type: got %q, want %q", got, want)
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
		return nil, errors.Wrapf(err, "could not resolve %q", src.GetTypeUrl())
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
