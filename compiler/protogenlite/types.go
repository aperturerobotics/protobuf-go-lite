// Package protogenlite provides a minimal protoc plugin model for generators
// that only need file/service/method metadata and Go import qualification.
package protogenlite

import "github.com/aperturerobotics/protobuf-go-lite/types/descriptorpb"

// GoImportPath identifies a Go package import path.
type GoImportPath string

// Ident constructs a Go identifier in the import path.
func (p GoImportPath) Ident(name string) GoIdent {
	return GoIdent{
		GoName:       name,
		GoImportPath: p,
	}
}

// GoIdent identifies a Go name within a package.
type GoIdent struct {
	GoName       string
	GoImportPath GoImportPath
}

// Comment stores a single comment string.
type Comment string

// String returns the comment text.
func (c Comment) String() string {
	return string(c)
}

// Comments stores descriptor comments.
type Comments struct {
	Leading Comment
}

// FileDesc exposes the minimal file descriptor helpers used by generators.
type FileDesc struct {
	path string
}

// Path returns the source proto path.
func (d *FileDesc) Path() string {
	if d == nil {
		return ""
	}
	return d.path
}

// ServiceDesc exposes the minimal service descriptor helpers used by generators.
type ServiceDesc struct {
	fullName string
}

// FullName returns the fully-qualified proto service name.
func (d *ServiceDesc) FullName() string {
	if d == nil {
		return ""
	}
	return d.fullName
}

// MethodDesc exposes the minimal method descriptor helpers used by generators.
type MethodDesc struct {
	name            string
	streamingClient bool
	streamingServer bool
}

// Name returns the proto method name.
func (d *MethodDesc) Name() string {
	if d == nil {
		return ""
	}
	return d.name
}

// IsStreamingClient reports whether the method accepts a client stream.
func (d *MethodDesc) IsStreamingClient() bool {
	if d == nil {
		return false
	}
	return d.streamingClient
}

// IsStreamingServer reports whether the method returns a server stream.
func (d *MethodDesc) IsStreamingServer() bool {
	if d == nil {
		return false
	}
	return d.streamingServer
}

// File stores the minimal generated view of a proto file.
type File struct {
	Proto                   *descriptorpb.FileDescriptorProto
	Desc                    *FileDesc
	Generate                bool
	GoImportPath            GoImportPath
	GoPackageName           string
	GeneratedFilenamePrefix string
	Messages                []*Message
	Services                []*Service
}

// Message stores the minimal generated view of a proto message.
type Message struct {
	Proto    *descriptorpb.DescriptorProto
	File     *File
	FullName string
	GoName   string
	GoIdent  GoIdent
	Messages []*Message
}

// Service stores the minimal generated view of a proto service.
type Service struct {
	Proto    *descriptorpb.ServiceDescriptorProto
	File     *File
	Desc     *ServiceDesc
	GoName   string
	Methods  []*Method
	Comments Comments
}

// Method stores the minimal generated view of a proto method.
type Method struct {
	Proto    *descriptorpb.MethodDescriptorProto
	Parent   *Service
	Desc     *MethodDesc
	GoName   string
	Input    *Message
	Output   *Message
	Comments Comments
}
