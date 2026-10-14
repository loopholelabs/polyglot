// SPDX-License-Identifier: Apache-2.0

package typescript

import (
	"bytes"
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
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name:    proto.String("oneof.proto"),
			Package: proto.String("test"),
			Syntax:  proto.String("proto3"),
			Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/test;test")},
			EnumType: []*descriptorpb.EnumDescriptorProto{{
				Name: proto.String("Test"),
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: proto.String("ZERO"), Number: proto.Int32(0)},
					{Name: proto.String("ONE"), Number: proto.Int32(1)},
				},
			}},
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("Data"),
					Field: []*descriptorpb.FieldDescriptorProto{
						tsField("message", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, nil, ""),
					},
				},
				{
					Name: proto.String("Choice"),
					OneofDecl: []*descriptorpb.OneofDescriptorProto{
						{Name: proto.String("value")},
						{Name: proto.String("_label")},
					},
					Field: []*descriptorpb.FieldDescriptorProto{
						tsField("ordinary", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, nil, ""),
						tsField("text", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING, &valueOneof, ""),
						tsField("raw", 7, descriptorpb.FieldDescriptorProto_TYPE_BYTES, &valueOneof, ""),
						tsField("payload", 11, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, &valueOneof, ".test.Data"),
						tsField("enum_value", 15, descriptorpb.FieldDescriptorProto_TYPE_ENUM, &valueOneof, ".test.Test"),
						tsField("enabled", 20, descriptorpb.FieldDescriptorProto_TYPE_BOOL, &valueOneof, ""),
						tsField("number", 25, descriptorpb.FieldDescriptorProto_TYPE_INT64, &valueOneof, ""),
						{
							Name:           proto.String("label"),
							Number:         proto.Int32(40),
							Label:          descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:           descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							OneofIndex:     &optionalOneof,
							Proto3Optional: proto.Bool(true),
						},
					},
				},
			},
		}},
	}

	gen := New()
	plugin, err := gen.options.New(req)
	require.NoError(t, err)
	require.Len(t, plugin.Files, 1)
	file := plugin.Files[0]
	var body bytes.Buffer
	err = gen.templ.ExecuteTemplate(&body, "body.templ", map[string]interface{}{
		"enums":    file.Desc.Enums(),
		"messages": file.Desc.Messages(),
	})
	require.NoError(t, err)
	generated := body.String()

	require.Contains(t, generated, "export enum TestChoiceValueCase")
	require.Contains(t, generated, "Text = 2")
	require.Contains(t, generated, "Raw = 7")
	require.Contains(t, generated, "private _polyglotValueCase")
	require.Contains(t, generated, "encodeUint32(encoded, this._polyglotValueCase)")
	require.Contains(t, generated, "setText(value: string | undefined)")
	require.Contains(t, generated, "switch (valueCase)")
	require.NotContains(t, generated, "TestChoiceLabelCase")
}

func tsField(name string, number int32, kind descriptorpb.FieldDescriptorProto_Type, oneof *int32, typeName string) *descriptorpb.FieldDescriptorProto {
	f := &descriptorpb.FieldDescriptorProto{
		Name:       proto.String(name),
		Number:     proto.Int32(number),
		Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:       kind.Enum(),
		OneofIndex: oneof,
	}
	if typeName != "" {
		f.TypeName = proto.String(typeName)
	}
	return f
}
