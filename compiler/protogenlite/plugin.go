package protogenlite

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/aperturerobotics/protobuf-go-lite/types/descriptorpb"
	pluginpb "github.com/aperturerobotics/protobuf-go-lite/types/pluginpb"
)

// Options configures a plugin run.
type Options struct{}

// Run executes a protoc plugin using the minimal descriptor model.
func (opts Options) Run(f func(*Plugin) error) {
	if err := run(opts, f); err != nil {
		_, _ = os.Stderr.WriteString(filepathBase(os.Args[0]) + ": " + err.Error() + "\n")
		os.Exit(1)
	}
}

func run(opts Options, f func(*Plugin) error) error {
	if len(os.Args) > 1 {
		return fmt.Errorf("unknown argument %q (this program should be run by protoc, not directly)", os.Args[1])
	}
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := req.UnmarshalVT(in); err != nil {
		return err
	}
	plugin, err := newPlugin(req)
	if err != nil {
		return err
	}
	if err := f(plugin); err != nil {
		msg := err.Error()
		plugin.err = &msg
	}
	resp := plugin.Response()
	out, err := resp.MarshalVT()
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

// Plugin stores a protoc plugin invocation.
type Plugin struct {
	Request           *pluginpb.CodeGeneratorRequest
	Files             []*File
	SupportedFeatures uint64

	err            *string
	packageNames   map[GoImportPath]string
	generatedFiles []*generatedFileRef
}

type generatedFileRef struct {
	filename string
	file     *GeneratedFile
}

func newPlugin(req *pluginpb.CodeGeneratorRequest) (*Plugin, error) {
	plugin := &Plugin{
		Request:      req,
		packageNames: make(map[GoImportPath]string),
	}

	generateSet := make(map[string]struct{}, len(req.GetFileToGenerate()))
	for _, name := range req.GetFileToGenerate() {
		generateSet[name] = struct{}{}
	}

	filesByName := make(map[string]*File, len(req.GetProtoFile()))
	for _, fileProto := range req.GetProtoFile() {
		goImportPath, goPackageName := deriveGoPackage(fileProto)
		file := &File{
			Proto:                   fileProto,
			Desc:                    &FileDesc{path: fileProto.GetName()},
			GoImportPath:            goImportPath,
			GoPackageName:           goPackageName,
			GeneratedFilenamePrefix: strings.TrimSuffix(fileProto.GetName(), ".proto"),
		}
		_, file.Generate = generateSet[fileProto.GetName()]
		filesByName[fileProto.GetName()] = file
		plugin.Files = append(plugin.Files, file)
		plugin.packageNames[file.GoImportPath] = file.GoPackageName
	}

	messagesByFullName := make(map[string]*Message)
	for _, file := range plugin.Files {
		file.Messages = buildMessages(file, file.Proto.GetMessageType(), nil, file.Proto.GetPackage(), messagesByFullName)
	}

	for _, file := range plugin.Files {
		file.Services = buildServices(file, messagesByFullName)
	}

	return plugin, nil
}

// NewGeneratedFile allocates a generated file in the plugin response.
func (p *Plugin) NewGeneratedFile(filename string, importPath GoImportPath) *GeneratedFile {
	file := newGeneratedFile(p, filename, importPath)
	p.generatedFiles = append(p.generatedFiles, &generatedFileRef{
		filename: filename,
		file:     file,
	})
	return file
}

// Response builds the code generator response.
func (p *Plugin) Response() *pluginpb.CodeGeneratorResponse {
	resp := &pluginpb.CodeGeneratorResponse{}
	if p.err != nil {
		resp.Error = p.err
		return resp
	}
	if p.SupportedFeatures != 0 {
		supportedFeatures := p.SupportedFeatures
		resp.SupportedFeatures = &supportedFeatures
	}
	for _, generated := range p.generatedFiles {
		content, err := generated.file.Content()
		if err != nil {
			msg := err.Error()
			return &pluginpb.CodeGeneratorResponse{Error: &msg}
		}
		name := generated.filename
		text := string(content)
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    &name,
			Content: &text,
		})
	}
	return resp
}

func (p *Plugin) packageNameFor(importPath GoImportPath) string {
	if pkg := p.packageNames[importPath]; pkg != "" {
		return pkg
	}
	return goPackageName(string(importPath))
}

func buildMessages(
	file *File,
	protos []*descriptorpb.DescriptorProto,
	parent *Message,
	protoPackage string,
	messagesByFullName map[string]*Message,
) []*Message {
	messages := make([]*Message, 0, len(protos))
	for _, proto := range protos {
		goName := goCamelCase(proto.GetName())
		if parent != nil {
			goName = parent.GoName + "_" + goName
		}
		fullName := proto.GetName()
		if parent != nil {
			fullName = parent.FullName + "." + proto.GetName()
		} else if protoPackage != "" {
			fullName = protoPackage + "." + proto.GetName()
		}
		message := &Message{
			Proto:    proto,
			File:     file,
			FullName: fullName,
			GoName:   goName,
			GoIdent: GoIdent{
				GoName:       goName,
				GoImportPath: file.GoImportPath,
			},
		}
		message.Messages = buildMessages(file, proto.GetNestedType(), message, protoPackage, messagesByFullName)
		messages = append(messages, message)
		messagesByFullName[fullName] = message
	}
	return messages
}

func buildServices(file *File, messagesByFullName map[string]*Message) []*Service {
	protos := file.Proto.GetService()
	services := make([]*Service, 0, len(protos))
	for i, proto := range protos {
		fullName := proto.GetName()
		if pkg := file.Proto.GetPackage(); pkg != "" {
			fullName = pkg + "." + fullName
		}
		service := &Service{
			Proto:    proto,
			File:     file,
			Desc:     &ServiceDesc{fullName: fullName},
			GoName:   goCamelCase(proto.GetName()),
			Comments: Comments{Leading: lookupLeadingComment(file.Proto, []int32{6, int32(i)})},
		}
		for j, methodProto := range proto.GetMethod() {
			input := messagesByFullName[strings.TrimPrefix(methodProto.GetInputType(), ".")]
			output := messagesByFullName[strings.TrimPrefix(methodProto.GetOutputType(), ".")]
			if input == nil || output == nil {
				continue
			}
			service.Methods = append(service.Methods, &Method{
				Proto:  methodProto,
				Parent: service,
				Desc: &MethodDesc{
					name:            methodProto.GetName(),
					streamingClient: methodProto.GetClientStreaming(),
					streamingServer: methodProto.GetServerStreaming(),
				},
				GoName:   goCamelCase(methodProto.GetName()),
				Input:    input,
				Output:   output,
				Comments: Comments{Leading: lookupLeadingComment(file.Proto, []int32{6, int32(i), 2, int32(j)})},
			})
		}
		services = append(services, service)
	}
	return services
}

func deriveGoPackage(file *descriptorpb.FileDescriptorProto) (GoImportPath, string) {
	if importPath, packageName, ok := wellKnownType(file.GetName()); ok {
		return GoImportPath(importPath), packageName
	}
	if option := file.GetOptions().GetGoPackage(); option != "" {
		if importPath, packageName, ok := strings.Cut(option, ";"); ok {
			return GoImportPath(importPath), packageName
		}
		importPath := option
		return GoImportPath(importPath), goPackageName(importPath)
	}
	importPath := path.Dir(file.GetName())
	if importPath == "." {
		importPath = ""
	}
	packageName := file.GetPackage()
	if packageName == "" {
		packageName = path.Base(importPath)
	}
	return GoImportPath(importPath), goPackageName(strings.ReplaceAll(packageName, ".", "_"))
}

func wellKnownType(fileName string) (string, string, bool) {
	switch fileName {
	case "google/protobuf/any.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/anypb", "anypb", true
	case "google/protobuf/api.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/apipb", "apipb", true
	case "google/protobuf/duration.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/durationpb", "durationpb", true
	case "google/protobuf/empty.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/emptypb", "emptypb", true
	case "google/protobuf/source_context.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/sourcecontextpb", "sourcecontextpb", true
	case "google/protobuf/struct.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/structpb", "structpb", true
	case "google/protobuf/timestamp.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb", "timestamppb", true
	case "google/protobuf/type.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/typepb", "typepb", true
	case "google/protobuf/wrappers.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/known/wrapperspb", "wrapperspb", true
	case "google/protobuf/descriptor.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/descriptorpb", "descriptorpb", true
	case "google/protobuf/compiler/plugin.proto":
		return "github.com/aperturerobotics/protobuf-go-lite/types/pluginpb", "google_protobuf_compiler", true
	default:
		return "", "", false
	}
}

func lookupLeadingComment(file *descriptorpb.FileDescriptorProto, pathValues []int32) Comment {
	sourceInfo := file.GetSourceCodeInfo()
	for _, location := range sourceInfo.GetLocation() {
		if samePath(location.GetPath(), pathValues) {
			return Comment(location.GetLeadingComments())
		}
	}
	return ""
}

func samePath(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func goPackageName(s string) string {
	if s == "" {
		return ""
	}
	s = path.Base(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)
	r, _ := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(r) {
		s = "_" + s
	}
	return s
}

func goCamelCase(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
		case c == '.':
			out = append(out, '_')
		case c == '_' && (i == 0 || s[i-1] == '.'):
			out = append(out, 'X')
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
		case isASCIIDigit(c):
			out = append(out, c)
		default:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			out = append(out, c)
			for i+1 < len(s) && isASCIILower(s[i+1]) {
				i++
				out = append(out, s[i])
			}
		}
	}
	return string(out)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func filepathBase(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return s[i+1:]
		}
	}
	return s
}
