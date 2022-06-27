/*
	Copyright 2022 Loophole Labs

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

package polyglot

type encoder C

func Encoder(b *C) *encoder {
	return (*encoder)(b)
}

func (e *encoder) Nil() *encoder {
	encodeNil((*C)(e))
	return e
}

func (e *encoder) Map(size uint32, keyKind, valueKind Kind) *encoder {
	encodeMap((*C)(e), size, keyKind, valueKind)
	return e
}

func (e *encoder) Slice(size uint32, kind Kind) *encoder {
	encodeSlice((*C)(e), size, kind)
	return e
}

func (e *encoder) Bytes(value []byte) *encoder {
	encodeBytes((*C)(e), value)
	return e
}

func (e *encoder) String(value string) *encoder {
	encodeString((*C)(e), value)
	return e
}

func (e *encoder) Error(value error) *encoder {
	encodeError((*C)(e), value)
	return e
}

func (e *encoder) Bool(value bool) *encoder {
	encodeBool((*C)(e), value)
	return e
}

func (e *encoder) Uint8(value uint8) *encoder {
	encodeUint8((*C)(e), value)
	return e
}

func (e *encoder) Uint16(value uint16) *encoder {
	encodeUint16((*C)(e), value)
	return e
}

func (e *encoder) Uint32(value uint32) *encoder {
	encodeUint32((*C)(e), value)
	return e
}

func (e *encoder) Uint64(value uint64) *encoder {
	encodeUint64((*C)(e), value)
	return e
}

func (e *encoder) Int32(value int32) *encoder {
	encodeInt32((*C)(e), value)
	return e
}

func (e *encoder) Int64(value int64) *encoder {
	encodeInt64((*C)(e), value)
	return e
}

func (e *encoder) Float32(value float32) *encoder {
	encodeFloat32((*C)(e), value)
	return e
}

func (e *encoder) Float64(value float64) *encoder {
	encodeFloat64((*C)(e), value)
	return e
}
