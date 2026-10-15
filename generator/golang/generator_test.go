// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerateOneof(t *testing.T) {
	valueOneof := int32(0)
	optionalOneof := int32(1)
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"oneof.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			{
				Name:    proto.String("oneof.proto"),
				Package: proto.String("test"),
				Syntax:  proto.String("proto3"),
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("example.com/test;test"),
				},
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Choice"),
						OneofDecl: []*descriptorpb.OneofDescriptorProto{
							{Name: proto.String("value")},
							{Name: proto.String("_label")},
						},
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:       proto.String("text"),
								Number:     proto.Int32(2),
								Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:       descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								OneofIndex: &valueOneof,
							},
							{
								Name:       proto.String("error"),
								Number:     proto.Int32(9),
								Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:       descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								OneofIndex: &valueOneof,
							},
							{
								Name:           proto.String("label"),
								Number:         proto.Int32(12),
								Label:          descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:           descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								OneofIndex:     &optionalOneof,
								Proto3Optional: proto.Bool(true),
							},
							{
								Name:   proto.String("value_case"),
								Number: proto.Int32(13),
								Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							},
						},
					},
				},
			},
		},
	}

	gen := New()
	plugin, err := gen.options.New(req)
	require.NoError(t, err)
	require.Len(t, plugin.Files, 1)
	file := plugin.Files[0]
	require.True(t, file.Generate)

	var out bytes.Buffer
	err = gen.templ.ExecuteTemplate(&out, "base.templ", map[string]interface{}{
		"pluginVersion":   "test",
		"sourcePath":      file.Desc.Path(),
		"package":         "test",
		"requiredImports": RequiredImports,
		"enums":           file.Desc.Enums(),
		"messages":        file.Desc.Messages(),
		"header":          true,
	})
	require.NoError(t, err)
	generated := out.String()
	_, err = parser.ParseFile(token.NewFileSet(), "oneof.polyglot.go", generated, 0)
	require.NoError(t, err)

	require.Contains(t, generated, "type TestChoiceValueCase uint32")
	require.Contains(t, generated, "TestChoiceValueCaseText TestChoiceValueCase = 2")
	require.Contains(t, generated, "TestChoiceValueCaseError TestChoiceValueCase = 9")
	require.Contains(t, generated, "valueCase TestChoiceValueCase")
	require.Contains(t, generated, "ValueCase string")
	require.Contains(t, generated, "switch x.valueCase")
	require.Contains(t, generated, "func (x *TestChoice) SetText(value string)")
	require.Contains(t, generated, "func (x *TestChoice) GetError() string")
	require.Contains(t, generated, "func (x *TestChoice) encodeError(")
	require.NotContains(t, generated, "func (x *TestChoice) Error(")
	require.NotContains(t, generated, "TestChoiceLabelCase")
}
