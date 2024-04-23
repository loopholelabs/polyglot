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

package rust

import (
	"github.com/loopholelabs/polyglot/v2/utils"

	"google.golang.org/protobuf/reflect/protoreflect"

	"errors"
	"fmt"
)

var (
	errUnknownKind        = errors.New("unknown or unsupported protoreflect.Kind")
	errUnknownCardinality = errors.New("unknown or unsupported protoreflect.Cardinality")
)

var (
	typeLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     "bool",
		protoreflect.Int32Kind:    "i32",
		protoreflect.Sint32Kind:   "i32",
		protoreflect.Uint32Kind:   "u32",
		protoreflect.Int64Kind:    "i64",
		protoreflect.Sint64Kind:   "i64",
		protoreflect.Uint64Kind:   "u64",
		protoreflect.Sfixed32Kind: "i32",
		protoreflect.Sfixed64Kind: "i64",
		protoreflect.Fixed32Kind:  "u32",
		protoreflect.Fixed64Kind:  "u64",
		protoreflect.FloatKind:    "f32",
		protoreflect.DoubleKind:   "f64",
		protoreflect.StringKind:   "String",
		protoreflect.BytesKind:    "Vec<u8>",
	}

	encodeLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     ".encode_bool",
		protoreflect.Int32Kind:    ".encode_i32",
		protoreflect.Sint32Kind:   ".encode_i32",
		protoreflect.Uint32Kind:   ".encode_u32",
		protoreflect.Int64Kind:    ".encode_i64",
		protoreflect.Sint64Kind:   ".encode_i64",
		protoreflect.Uint64Kind:   ".encode_u64",
		protoreflect.Sfixed32Kind: ".encode_i32",
		protoreflect.Sfixed64Kind: ".encode_i64",
		protoreflect.Fixed32Kind:  ".encode_u32",
		protoreflect.Fixed64Kind:  ".encode_u64",
		protoreflect.StringKind:   ".encode_string",
		protoreflect.FloatKind:    ".encode_f32",
		protoreflect.DoubleKind:   ".encode_f64",
		protoreflect.BytesKind:    ".encode_bytes",
		protoreflect.EnumKind:     ".encode_u32",
	}

	decodeLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     ".decode_bool",
		protoreflect.Int32Kind:    ".decode_i32",
		protoreflect.Sint32Kind:   ".decode_i32",
		protoreflect.Uint32Kind:   ".decode_u32",
		protoreflect.Int64Kind:    ".decode_i64",
		protoreflect.Sint64Kind:   ".decode_i64",
		protoreflect.Uint64Kind:   ".decode_u64",
		protoreflect.Sfixed32Kind: ".decode_i32",
		protoreflect.Sfixed64Kind: ".decode_i64",
		protoreflect.Fixed32Kind:  ".decode_u32",
		protoreflect.Fixed64Kind:  ".decode_u64",
		protoreflect.StringKind:   ".decode_string",
		protoreflect.FloatKind:    ".decode_f32",
		protoreflect.DoubleKind:   ".decode_f64",
		protoreflect.BytesKind:    ".decode_bytes",
		protoreflect.EnumKind:     ".decode_u32",
	}

	kindLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     "Kind::Bool",
		protoreflect.Int32Kind:    "Kind::I32",
		protoreflect.Sint32Kind:   "Kind::I32",
		protoreflect.Uint32Kind:   "Kind::U32",
		protoreflect.Int64Kind:    "Kind::I64",
		protoreflect.Sint64Kind:   "Kind::I64",
		protoreflect.Uint64Kind:   "Kind::U64",
		protoreflect.Sfixed32Kind: "Kind::I32",
		protoreflect.Sfixed64Kind: "Kind::I64",
		protoreflect.Fixed32Kind:  "Kind::U32",
		protoreflect.Fixed64Kind:  "Kind::U64",
		protoreflect.StringKind:   "Kind::String",
		protoreflect.FloatKind:    "Kind::F32",
		protoreflect.DoubleKind:   "Kind::F64",
		protoreflect.BytesKind:    "Kind::Bytes",
		protoreflect.EnumKind:     "Kind::U32",
	}
)

func findValue(field protoreflect.FieldDescriptor) string {
	if kind, ok := typeLUT[field.Kind()]; !ok {
		switch field.Kind() {
		case protoreflect.EnumKind:
			switch field.Cardinality() {
			case protoreflect.Optional, protoreflect.Required:
				return utils.CamelCase(string(field.Enum().FullName()))
			case protoreflect.Repeated:
				return fmt.Sprintf("Vec<%s>", utils.CamelCase(string(field.Enum().FullName())))
			default:
				panic(errUnknownCardinality)
			}
		case protoreflect.MessageKind:
			if field.IsMap() {
				return fmt.Sprintf("HashMap<%s, %s>", findValue(field.MapKey()), findValue(field.MapValue()))
			} else {
				switch field.Cardinality() {
				case protoreflect.Optional, protoreflect.Required:
					return utils.CamelCase(string(field.Message().FullName()))
				case protoreflect.Repeated:
					return fmt.Sprintf("Vec<%s>", utils.CamelCase(string(field.Message().FullName())))
				default:
					panic(errUnknownCardinality)
				}
			}
		default:
			panic(errUnknownKind)
		}
	} else {
		if field.Cardinality() == protoreflect.Repeated {
			kind = "Vec<" + kind + ">"
		}
		return kind
	}
}

type encodingFields struct {
	MessageFields []protoreflect.FieldDescriptor
	SliceFields   []protoreflect.FieldDescriptor
	Values        []string
}

func getEncodingFields(fields protoreflect.FieldDescriptors) encodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var values []string

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else {
			if encoder, ok := encodeLUT[field.Kind()]; !ok {
				switch field.Kind() {
				case protoreflect.MessageKind:
					messageFields = append(messageFields, field)
				default:
					panic(errUnknownKind)
				}
			} else {
				if field.Kind() == protoreflect.EnumKind {
					values = append(values, fmt.Sprintf("%s(self.%s  as u32)", encoder, utils.SnakeCaseName(field.Name())))
				} else if field.Kind() == protoreflect.StringKind {
					values = append(values, fmt.Sprintf("%s(&self.%s)", encoder, utils.SnakeCaseName(field.Name())))
				} else if field.Kind() == protoreflect.BytesKind {
					values = append(values, fmt.Sprintf("%s(&self.%s)", encoder, utils.SnakeCaseName(field.Name())))
				} else {
					values = append(values, fmt.Sprintf("%s(self.%s)", encoder, utils.SnakeCaseName(field.Name())))
				}
			}
		}
	}
	return encodingFields{
		MessageFields: messageFields,
		SliceFields:   sliceFields,
		Values:        values,
	}
}

type decodingFields struct {
	MessageFields []protoreflect.FieldDescriptor
	SliceFields   []protoreflect.FieldDescriptor
	Other         []protoreflect.FieldDescriptor
}

func getDecodingFields(fields protoreflect.FieldDescriptors) decodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var other []protoreflect.FieldDescriptor

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else {
			if _, ok := decodeLUT[field.Kind()]; !ok {
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

	return decodingFields{
		MessageFields: messageFields,
		SliceFields:   sliceFields,
		Other:         other,
	}
}

func getKind(kind protoreflect.Kind) string {
	var outKind string
	var ok bool
	if outKind, ok = kindLUT[kind]; !ok {
		switch kind {
		case protoreflect.MessageKind:
			outKind = polyglotAnyKind
		default:
			panic(errUnknownKind)
		}
	}
	return outKind
}

func getLUTEncoder(kind protoreflect.Kind) string {
	return encodeLUT[kind]
}

func getLUTDecoder(kind protoreflect.Kind) string {
	return decodeLUT[kind]
}

func getKindLUT(kind protoreflect.Kind) string {
	return kindLUT[kind]
}
