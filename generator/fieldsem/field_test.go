package fieldsem

import (
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type testQualifier struct{}

func (testQualifier) QualifiedGoIdent(ident protogen.GoIdent) string {
	return ident.GoName
}

func TestResolveEditionFieldSemantics(t *testing.T) {
	fields := testEditionFields(t)
	tests := []struct {
		name string
		want Field
	}{
		{
			name: "explicit_int32",
			want: Field{Type: "int32", Pointer: true, Reference: true, EmitDefault: true},
		},
		{
			name: "implicit_int32",
			want: Field{Type: "int32"},
		},
		{
			name: "required_int32",
			want: Field{Type: "int32", Pointer: true, Reference: true, Required: true, EmitDefault: true},
		},
		{
			name: "explicit_bytes",
			want: Field{Type: "[]byte", Reference: true, EmitDefault: true},
		},
		{
			name: "nested_message",
			want: Field{Type: "*Msg_Nested", Reference: true, EmitDefault: true},
		},
		{
			name: "packed_int32",
			want: Field{Type: "[]int32", Reference: true, Packed: true, List: true},
		},
		{
			name: "expanded_int32",
			want: Field{Type: "[]int32", Reference: true, List: true},
		},
		{
			name: "nested_map",
			want: Field{Type: "map[string]*Msg_Nested", Reference: true, Map: true},
		},
		{
			name: "choice_int32",
			want: Field{Type: "int32", Reference: true, RealOneof: true, EmitDefault: true},
		},
	}

	for _, test := range tests {
		field := fields[test.name]
		if field == nil {
			t.Fatalf("missing field %s", test.name)
		}
		got := Resolve(testQualifier{}, field)
		if got != test.want {
			t.Errorf("Resolve(%s) = %+v, want %+v", test.name, got, test.want)
		}
	}
}

func testEditionFields(t *testing.T) map[string]*protogen.Field {
	t.Helper()

	file := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("semantics.proto"),
		Syntax:  proto.String("editions"),
		Package: proto.String("semantics"),
		Edition: descriptorpb.Edition_EDITION_2024.Enum(),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/semantics"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("Msg"),
				Field: []*descriptorpb.FieldDescriptorProto{
					scalarField("explicit_int32", 1, descriptorpb.FieldDescriptorProto_TYPE_INT32, nil),
					scalarField("implicit_int32", 2, descriptorpb.FieldDescriptorProto_TYPE_INT32, &descriptorpb.FeatureSet{
						FieldPresence: descriptorpb.FeatureSet_IMPLICIT.Enum(),
					}),
					scalarField("required_int32", 3, descriptorpb.FieldDescriptorProto_TYPE_INT32, &descriptorpb.FeatureSet{
						FieldPresence: descriptorpb.FeatureSet_LEGACY_REQUIRED.Enum(),
					}),
					scalarField("explicit_bytes", 4, descriptorpb.FieldDescriptorProto_TYPE_BYTES, nil),
					{
						Name:     proto.String("nested_message"),
						Number:   proto.Int32(5),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".semantics.Msg.Nested"),
						JsonName: proto.String("nestedMessage"),
					},
					repeatedField("packed_int32", 6, nil),
					repeatedField("expanded_int32", 7, &descriptorpb.FeatureSet{
						RepeatedFieldEncoding: descriptorpb.FeatureSet_EXPANDED.Enum(),
					}),
					{
						Name:     proto.String("nested_map"),
						Number:   proto.Int32(8),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".semantics.Msg.NestedMapEntry"),
						JsonName: proto.String("nestedMap"),
					},
					{
						Name:       proto.String("choice_int32"),
						Number:     proto.Int32(9),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
						OneofIndex: proto.Int32(0),
						JsonName:   proto.String("choiceInt32"),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Nested"),
						Field: []*descriptorpb.FieldDescriptorProto{
							scalarField("value", 1, descriptorpb.FieldDescriptorProto_TYPE_INT32, nil),
						},
					},
					{
						Name: proto.String("NestedMapEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							scalarField("key", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, nil),
							{
								Name:     proto.String("value"),
								Number:   proto.Int32(2),
								Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
								TypeName: proto.String(".semantics.Msg.Nested"),
								JsonName: proto.String("value"),
							},
						},
						Options: &descriptorpb.MessageOptions{
							MapEntry: proto.Bool(true),
						},
					},
				},
				OneofDecl: []*descriptorpb.OneofDescriptorProto{
					{Name: proto.String("choice")},
				},
			},
		},
	}
	plugin, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
		ProtoFile:      []*descriptorpb.FileDescriptorProto{file},
		FileToGenerate: []string{"semantics.proto"},
	})
	if err != nil {
		t.Fatal(err)
	}

	fields := make(map[string]*protogen.Field)
	for _, field := range plugin.Files[0].Messages[0].Fields {
		fields[string(field.Desc.Name())] = field
	}
	return fields
}

func scalarField(
	name string,
	num int32,
	typ descriptorpb.FieldDescriptorProto_Type,
	features *descriptorpb.FeatureSet,
) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Number:   proto.Int32(num),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:     typ.Enum(),
		JsonName: proto.String(name),
		Options: &descriptorpb.FieldOptions{
			Features: features,
		},
	}
}

func repeatedField(name string, num int32, features *descriptorpb.FeatureSet) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Number:   proto.Int32(num),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
		JsonName: proto.String(name),
		Options: &descriptorpb.FieldOptions{
			Features: features,
		},
	}
}
