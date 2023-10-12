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

package polyglot

import (
	"math"
	"reflect"
	"unsafe"
)

var (
	NilRawKind     = byte(0)
	SliceRawKind   = byte(1)
	MapRawKind     = byte(2)
	AnyRawKind     = byte(3)
	BytesRawKind   = byte(4)
	StringRawKind  = byte(5)
	ErrorRawKind   = byte(6)
	BoolRawKind    = byte(7)
	Uint8RawKind   = byte(8)
	Uint16RawKind  = byte(9)
	Uint32RawKind  = byte(10)
	Uint64RawKind  = byte(11)
	Int32RawKind   = byte(12)
	Int64RawKind   = byte(13)
	Float32RawKind = byte(14)
	Float64RawKind = byte(15)
)

type Kind byte

var (
	NilKind     = Kind(NilRawKind)
	SliceKind   = Kind(SliceRawKind)
	MapKind     = Kind(MapRawKind)
	AnyKind     = Kind(AnyRawKind)
	BytesKind   = Kind(BytesRawKind)
	StringKind  = Kind(StringRawKind)
	ErrorKind   = Kind(ErrorRawKind)
	BoolKind    = Kind(BoolRawKind)
	Uint8Kind   = Kind(Uint8RawKind)
	Uint16Kind  = Kind(Uint16RawKind)
	Uint32Kind  = Kind(Uint32RawKind)
	Uint64Kind  = Kind(Uint64RawKind)
	Int32Kind   = Kind(Int32RawKind)
	Int64Kind   = Kind(Int64RawKind)
	Float32Kind = Kind(Float32RawKind)
	Float64Kind = Kind(Float64RawKind)
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) Is(err error) bool {
	return e.Error() == err.Error()
}

var (
	falseBool = byte(0)
	trueBool  = byte(1)
)

func encodeNil(b *Buffer) {
	b.Grow(1)
	b.WriteRawByte(NilRawKind)
}

func encodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	b.Grow(3)
	b.WriteRawByte(MapRawKind)
	b.WriteRawByte(byte(keyKind))
	b.WriteRawByte(byte(valueKind))
	encodeUint32(b, size)
}

func encodeSlice(b *Buffer, size uint32, kind Kind) {
	b.Grow(2)
	b.WriteRawByte(SliceRawKind)
	b.WriteRawByte(byte(kind))
	encodeUint32(b, size)
}

func encodeBytes(b *Buffer, value []byte) {
	b.Grow(1)
	b.WriteRawByte(BytesRawKind)
	encodeUint32(b, uint32(len(value)))
	b.Write(value)
}

func encodeString(b *Buffer, value string) {
	var nb []byte
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&nb))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&value))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	b.Grow(1)
	b.WriteRawByte(StringRawKind)
	encodeUint32(b, uint32(len(nb)))
	b.Write(nb)
}

func encodeError(b *Buffer, err error) {
	b.Grow(1)
	b.WriteRawByte(ErrorRawKind)
	encodeString(b, err.Error())
}

func encodeBool(b *Buffer, value bool) {
	b.Grow(2)
	b.WriteRawByte(BoolRawKind)
	if value {
		b.WriteRawByte(trueBool)
	} else {
		b.WriteRawByte(falseBool)
	}
}

func encodeUint8(b *Buffer, value uint8) {
	b.Grow(2)
	b.WriteRawByte(Uint8RawKind)
	b.WriteRawByte(value)
}

// Variable integer encoding with the same format as binary.varint
// (https://developers.google.com/protocol-buffers/docs/encoding#varints)
func encodeUint16(b *Buffer, value uint16) {
	b.Grow(VarIntLen16)
	b.WriteRawByte(Uint16RawKind)
	for value >= continuation {
		// Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
		b.WriteRawByte(byte(value) | continuation)
		value >>= 7
	}
	b.WriteRawByte(byte(value))
}

func encodeUint32(b *Buffer, value uint32) {
	b.Grow(VarIntLen32)
	b.WriteRawByte(Uint32RawKind)
	for value >= continuation {
		// Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
		b.WriteRawByte(byte(value) | continuation)
		value >>= 7
	}
	b.WriteRawByte(byte(value))
}

func encodeUint64(b *Buffer, value uint64) {
	b.Grow(VarIntLen64)
	b.WriteRawByte(Uint64RawKind)
	for value >= continuation {
		// Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
		b.WriteRawByte(byte(value) | continuation)
		value >>= 7
	}
	b.WriteRawByte(byte(value))
}

func encodeInt32(b *Buffer, value int32) {
	b.Grow(VarIntLen32)
	b.WriteRawByte(Int32RawKind)
	// Shift the value to the left by 1 bit, then flip the bits if the value is negative.
	castValue := uint32(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	for castValue >= continuation {
		// Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
		b.WriteRawByte(byte(castValue) | continuation)
		castValue >>= 7
	}
	b.WriteRawByte(byte(castValue))
}

func encodeInt64(b *Buffer, value int64) {
	b.Grow(VarIntLen64)
	b.WriteRawByte(Int64RawKind)
	// Shift the value to the left by 1 bit, then flip the bits if the value is negative.
	castValue := uint64(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	for castValue >= continuation {
		// Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
		b.WriteRawByte(byte(castValue) | continuation)
		castValue >>= 7
	}
	b.WriteRawByte(byte(castValue))
}

func encodeFloat32(b *Buffer, value float32) {
	b.Grow(5)
	b.WriteRawByte(Float32RawKind)
	castValue := math.Float32bits(value)
	b.WriteRawByte(byte(castValue >> 24))
	b.WriteRawByte(byte(castValue >> 16))
	b.WriteRawByte(byte(castValue >> 8))
	b.WriteRawByte(byte(castValue))
}

func encodeFloat64(b *Buffer, value float64) {
	b.Grow(9)
	b.WriteRawByte(Float64RawKind)
	castValue := math.Float64bits(value)
	b.WriteRawByte(byte(castValue >> 56))
	b.WriteRawByte(byte(castValue >> 48))
	b.WriteRawByte(byte(castValue >> 40))
	b.WriteRawByte(byte(castValue >> 32))
	b.WriteRawByte(byte(castValue >> 24))
	b.WriteRawByte(byte(castValue >> 16))
	b.WriteRawByte(byte(castValue >> 8))
	b.WriteRawByte(byte(castValue))
}
