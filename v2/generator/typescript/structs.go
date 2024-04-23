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

package typescript

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
		protoreflect.BoolKind:     "boolean",
		protoreflect.Int32Kind:    "number",
		protoreflect.Sint32Kind:   "number",
		protoreflect.Uint32Kind:   "number",
		protoreflect.Int64Kind:    "number",
		protoreflect.Sint64Kind:   "number",
		protoreflect.Uint64Kind:   "number",
		protoreflect.Sfixed32Kind: "number",
		protoreflect.Sfixed64Kind: "number",
		protoreflect.Fixed32Kind:  "number",
		protoreflect.Fixed64Kind:  "number",
		protoreflect.FloatKind:    "number",
		protoreflect.DoubleKind:   "number",
		protoreflect.StringKind:   "string",
		protoreflect.BytesKind:    "Uint8Array",
	}

	encodeLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     "encodeBoolean",
		protoreflect.Int32Kind:    "encodeInt32",
		protoreflect.Sint32Kind:   "encodeInt32",
		protoreflect.Uint32Kind:   "encodeUint32",
		protoreflect.Int64Kind:    "encodeInt64",
		protoreflect.Sint64Kind:   "encodeInt64",
		protoreflect.Uint64Kind:   "encodeUint64",
		protoreflect.Sfixed32Kind: "encodeInt32",
		protoreflect.Sfixed64Kind: "encodeInt64",
		protoreflect.Fixed32Kind:  "encodeUint32",
		protoreflect.Fixed64Kind:  "encodeUint64",
		protoreflect.StringKind:   "encodeString",
		protoreflect.FloatKind:    "encodeFloat32",
		protoreflect.DoubleKind:   "encodeFloat64",
		protoreflect.BytesKind:    "encodeBytes",
		protoreflect.EnumKind:     "encodeUint32",
	}

	decodeLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     "decodeBoolean",
		protoreflect.Int32Kind:    "decodeInt32",
		protoreflect.Sint32Kind:   "decodeInt32",
		protoreflect.Uint32Kind:   "decodeUint32",
		protoreflect.Int64Kind:    "decodeInt64",
		protoreflect.Sint64Kind:   "decodeInt64",
		protoreflect.Uint64Kind:   "decodeUint64",
		protoreflect.Sfixed32Kind: "decodeInt32",
		protoreflect.Sfixed64Kind: "decodeInt64",
		protoreflect.Fixed32Kind:  "decodeUint32",
		protoreflect.Fixed64Kind:  "decodeUint64",
		protoreflect.StringKind:   "decodeString",
		protoreflect.FloatKind:    "decodeFloat32",
		protoreflect.DoubleKind:   "decodeFloat64",
		protoreflect.BytesKind:    "decodeBytes",
		protoreflect.EnumKind:     "decodeUint32",
	}

	kindLUT = map[protoreflect.Kind]string{
		protoreflect.BoolKind:     "Kind.Bool",
		protoreflect.Int32Kind:    "Kind.Int32",
		protoreflect.Sint32Kind:   "Kind.Int32",
		protoreflect.Uint32Kind:   "Kind.U32",
		protoreflect.Int64Kind:    "Kind.I64",
		protoreflect.Sint64Kind:   "Kind.I64",
		protoreflect.Uint64Kind:   "Kind.U64",
		protoreflect.Sfixed32Kind: "Kind.Int32",
		protoreflect.Sfixed64Kind: "Kind.I64",
		protoreflect.Fixed32Kind:  "Kind.U32",
		protoreflect.Fixed64Kind:  "Kind.U64",
		protoreflect.StringKind:   "Kind.String",
		protoreflect.FloatKind:    "Kind.F32",
		protoreflect.DoubleKind:   "Kind.F64",
		protoreflect.BytesKind:    "Kind.Bytes",
		protoreflect.EnumKind:     "Kind.U32",
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
				return fmt.Sprintf("%s[]", utils.CamelCase(string(field.Enum().FullName())))
			default:
				panic(errUnknownCardinality)
			}
		case protoreflect.MessageKind:
			if field.IsMap() {
				return fmt.Sprintf("Map<%s, %s>", findValue(field.MapKey()), findValue(field.MapValue()))
			} else {
				switch field.Cardinality() {
				case protoreflect.Optional, protoreflect.Required:
					return utils.CamelCase(string(field.Message().FullName()))
				case protoreflect.Repeated:
					return fmt.Sprintf("%s[]", utils.CamelCase(string(field.Message().FullName())))
				default:
					panic(errUnknownCardinality)
				}
			}
		default:
			panic(errUnknownKind)
		}
	} else {
		if field.Cardinality() == protoreflect.Repeated {
			kind = kind + "[]"
		}
		return kind
	}
}

type encodingFields struct {
	MessageFields []protoreflect.FieldDescriptor
	SliceFields   []protoreflect.FieldDescriptor
	Values        []string
}

func getEncodingFields(trackDependency func(dep string) string, fields protoreflect.FieldDescriptors) encodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var values []string

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else {
			if encoder := getLUTEncoder(trackDependency, field.Kind()); encoder == "" {
				switch field.Kind() {
				case protoreflect.MessageKind:
					messageFields = append(messageFields, field)
				default:
					panic(errUnknownKind)
				}
			} else {
				if field.Kind() == protoreflect.EnumKind {
					values = append(values, fmt.Sprintf("%s(encoded, this.%s as number)", encoder, utils.LowercaseCamelCaseName(field.Name())))
				} else {
					values = append(values, fmt.Sprintf("%s(encoded, this.%s)", encoder, utils.LowercaseCamelCaseName(field.Name())))
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

func getDecodingFields(trackDependency func(dep string) string, fields protoreflect.FieldDescriptors) decodingFields {
	var messageFields []protoreflect.FieldDescriptor
	var sliceFields []protoreflect.FieldDescriptor
	var other []protoreflect.FieldDescriptor

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
			sliceFields = append(sliceFields, field)
		} else {
			if encoder := getLUTDecoder(trackDependency, field.Kind()); encoder == "" {
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

func getKind(dependencies map[string]struct{}, kind protoreflect.Kind) string {
	dependencies["Kind"] = struct{}{}

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

func getLUTEncoder(trackDependency func(dep string) string, kind protoreflect.Kind) string {
	encoder, ok := encodeLUT[kind]
	if ok {
		trackDependency(encoder)
	}

	return encoder
}

func getLUTDecoder(trackDependency func(dep string) string, kind protoreflect.Kind) string {
	decoder, ok := decodeLUT[kind]
	if ok {
		trackDependency(decoder)
	}

	return decoder
}

func getKindLUT(trackDependency func(dep string) string, kind protoreflect.Kind) string {
	trackDependency("Kind")

	return kindLUT[kind]
}
