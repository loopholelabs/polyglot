/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package golang

import (
	"errors"
	"fmt"

	gogen "github.com/loopholelabs/polyglot/generator/golang"
	"github.com/loopholelabs/polyglot/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	errUnknownKind        = errors.New("unknown or unsupported protoreflect.Kind")
	errUnknownCardinality = errors.New("unknown or unsupported protoreflect.Cardinality")
)

type EncodingFields struct {
	MessageFields []protoreflect.FieldDescriptor
	SliceFields   []protoreflect.FieldDescriptor
	MapFields     []protoreflect.FieldDescriptor
	Values        []string
}

func GetEncodingFields(g *Generator, fields protoreflect.FieldDescriptors) EncodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var mapFields []protoreflect.FieldDescriptor
	var values []string

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else if field.IsMap() {
			mapFields = append(mapFields, field)
		} else {
			if encoder, ok := gogen.EncodeLUT[field.Kind()]; !ok {
				switch field.Kind() {
				case protoreflect.MessageKind:
					messageFields = append(messageFields, field)
				default:
					panic(errUnknownKind)
				}
			} else {
				if field.Kind() == protoreflect.EnumKind {
					values = append(values, fmt.Sprintf("%s(uint32(x.%s))", encoder, g.camelCase(field.Name())))
				} else if field.Cardinality() == protoreflect.Optional && field.Kind() != protoreflect.BytesKind && *g.useNullablePointers {
					values = append(values, fmt.Sprintf("%sPtr(x.%s)", encoder, g.camelCase(field.Name())))
				} else {
					values = append(values, fmt.Sprintf("%s(x.%s)", encoder, g.camelCase(field.Name())))
				}
			}
		}
	}
	return EncodingFields{
		MessageFields: messageFields,
		SliceFields:   sliceFields,
		MapFields:     mapFields,
		Values:        values,
	}
}

type DecodingFields struct {
	MessageFields []protoreflect.FieldDescriptor
	SliceFields   []protoreflect.FieldDescriptor
	MapFields     []protoreflect.FieldDescriptor
	Other         []protoreflect.FieldDescriptor
}

func GetDecodingFields(fields protoreflect.FieldDescriptors) DecodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var mapFields []protoreflect.FieldDescriptor
	var other []protoreflect.FieldDescriptor

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else if field.IsMap() {
			mapFields = append(mapFields, field)
		} else {
			if _, ok := gogen.DecodeLUT[field.Kind()]; !ok {
				switch field.Kind() {
				case protoreflect.MessageKind:
					messageFields = append(messageFields, field)
				default:
					panic(errUnknownKind)
				}
			} else {
				other = append(other, field)
			}
		}
	}

	return DecodingFields{
		MessageFields: messageFields,
		SliceFields:   sliceFields,
		MapFields:     mapFields,
		Other:         other,
	}
}

func FindValue(g *Generator, field protoreflect.FieldDescriptor) string {
	if kind, ok := gogen.TypeLUT[field.Kind()]; !ok {
		switch field.Kind() {
		case protoreflect.EnumKind:
			switch field.Cardinality() {
			case protoreflect.Optional, protoreflect.Required:
				return g.protoName(field.Enum().FullName())
			case protoreflect.Repeated:
				return utils.CamelCase(utils.AppendString(Slice, g.protoName(field.Enum().FullName())))
			default:
				panic(errUnknownCardinality)
			}
		case protoreflect.MessageKind:
			if field.IsMap() {
				return utils.CamelCase(utils.AppendString(g.protoName(field.FullName()), MapSuffix))
			} else {
				switch field.Cardinality() {
				case protoreflect.Optional, protoreflect.Required:
					return utils.AppendString(Pointer, g.protoName(field.Message().FullName()))
				case protoreflect.Repeated:
					return utils.AppendString(Slice, Pointer, g.protoName(field.Message().FullName()))
				default:
					panic(errUnknownCardinality)
				}
			}
		default:
			panic(errUnknownKind)
		}
	} else {
		if field.Cardinality() == protoreflect.Repeated {
			kind = Slice + kind
		}
		return kind
	}
}
