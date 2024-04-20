package anypb_resolver

import (
	"errors"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
)

// ErrNoAnyTypeResolver is returned if no resolver was provided for the Any type.
var ErrNoAnyTypeResolver = errors.New("no resolver provided for the Any type")

// AnyTypeResolver is an interface for looking up messages.
//
// A compliant implementation must deterministically return the same type
// if no error is encountered.
//
// The [Types] type implements this interface.
type AnyTypeResolver interface {
	// FindMessageByURL looks up a message by a URL identifier.
	// See documentation on google.protobuf.Any.type_url for the URL format.
	//
	// Returns the constructor for the message.
	// This returns (nil, ErrNotFound) if not found.
	FindMessageByURL(url string) (func() protobuf_go_lite.Message, error)
}

// errAnyTypeResolver implements AnyTypeResolver returning an error.
type errAnyTypeResolver struct {
	err error
}

// FindMessageByURL looks up a message by a URL identifier.
func (e *errAnyTypeResolver) FindMessageByURL(url string) (func() protobuf_go_lite.Message, error) {
	return nil, e.err
}

// NewErrAnyTypeResolver constructs a new AnyTypeResolver that returns an error.
func NewErrAnyTypeResolver(err error) AnyTypeResolver {
	return &errAnyTypeResolver{err: err}
}

// funcAnyTypeResolver implements AnyTypeResolver with callbacks
type funcAnyTypeResolver struct {
	findMessageByURL func(url string) (func() protobuf_go_lite.Message, error)
}

// NewFuncAnyTypeResolver constructs a new AnyTypeResolver with callback funcs.
func NewFuncAnyTypeResolver(findMessageByURL func(url string) (func() protobuf_go_lite.Message, error)) AnyTypeResolver {
	return &funcAnyTypeResolver{findMessageByURL: findMessageByURL}
}

// FindMessageByURL looks up a message by a URL identifier.
func (e *funcAnyTypeResolver) FindMessageByURL(url string) (func() protobuf_go_lite.Message, error) {
	if e.findMessageByURL == nil {
		return nil, nil
	}
	return e.findMessageByURL(url)
}
