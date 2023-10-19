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

const (
	nilSize     = 1
	mapSize     = 3 + uint32Size
	sliceSize   = 2 + uint32Size
	bytesSize   = 1 + uint32Size
	stringSize  = 1 + uint32Size
	errorSize   = 1 + stringSize
	boolSize    = 2
	uint8Size   = 2
	uint16Size  = 1 + VarIntLen16
	uint32Size  = 1 + VarIntLen32
	uint64Size  = 1 + VarIntLen64
	float32Size = 5
	float64Size = 9
)

func EncodeNil(b *Buffer) {
	b.grow(nilSize)
	RawEncodeNil(b)
}

func RawEncodeNil(b *Buffer) {
	b.b[b.offset] = NilRawKind
	b.offset++
}

func EncodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	b.grow(mapSize)
	RawEncodeMap(b, size, keyKind, valueKind)
}

func RawEncodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	offset := b.offset
	b.b[offset] = MapRawKind
	offset++
	b.b[offset] = byte(keyKind)
	offset++
	b.b[offset] = byte(valueKind)
	b.offset = offset + 1
	RawEncodeUint32(b, size)
}

func EncodeSlice(b *Buffer, size uint32, kind Kind) {
	b.grow(sliceSize)
	RawEncodeSlice(b, size, kind)
}

func RawEncodeSlice(b *Buffer, size uint32, kind Kind) {
	offset := b.offset
	b.b[offset] = SliceRawKind
	offset++
	b.b[offset] = byte(kind)
	b.offset = offset + 1
	RawEncodeUint32(b, size)
}

func EncodeBytes(b *Buffer, value []byte) {
	b.grow(bytesSize + len(value))
	RawEncodeBytes(b, value)
}

func RawEncodeBytes(b *Buffer, value []byte) {
	b.b[b.offset] = BytesRawKind
	b.offset++
	RawEncodeUint32(b, uint32(len(value)))
	b.offset += copy(b.b[b.offset:], value)
}

func EncodeString(b *Buffer, value string) {
	b.grow(stringSize + len(value))
	RawEncodeString(b, value)
}

func RawEncodeString(b *Buffer, value string) {
	var nb []byte
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&nb))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&value))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	b.b[b.offset] = StringRawKind
	b.offset++
	RawEncodeUint32(b, uint32(len(nb)))
	b.offset += copy(b.b[b.offset:], nb)
}

func EncodeError(b *Buffer, err error) {
	errString := err.Error()
	b.grow(errorSize + len(errString))
	RawEncodeError(b, errString)
}

func RawEncodeError(b *Buffer, errString string) {
	b.b[b.offset] = ErrorRawKind
	b.offset++
	RawEncodeString(b, errString)
}

func EncodeBool(b *Buffer, value bool) {
	b.grow(boolSize)
	RawEncodeBool(b, value)
}

func RawEncodeBool(b *Buffer, value bool) {
	offset := b.offset
	b.b[offset] = BoolRawKind
	offset++
	if value {
		b.b[offset] = trueBool
	} else {
		b.b[offset] = falseBool
	}
	b.offset = offset + 1
}

func EncodeUint8(b *Buffer, value uint8) {
	b.grow(uint8Size)
	RawEncodeUint8(b, value)
}

func RawEncodeUint8(b *Buffer, value uint8) {
	offset := b.offset
	b.b[offset] = Uint8RawKind
	offset++
	b.b[offset] = value
	b.offset = offset + 1
}

func EncodeUint16(b *Buffer, value uint16) {
	b.grow(uint16Size)
	RawEncodeUint16(b, value)
}

func RawEncodeUint16(b *Buffer, value uint16) {
	offset := b.offset
	b.b[offset] = Uint16RawKind
	offset++
	if value < continuation {
		b.b[offset] = byte(value)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(value&(continuation-1)) | continuation
		offset++
		value >>= 7
		if value < continuation {
			b.b[offset] = byte(value)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(value&(continuation-1)) | continuation
			offset++
			value >>= 7
			if value < continuation {
				b.b[offset] = byte(value)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(value&(continuation-1)) | continuation
				offset++
				b.b[offset] = byte(value >> 7)
				b.offset = offset + 1
			}
		}
	}
}

func EncodeUint32(b *Buffer, value uint32) {
	b.grow(uint32Size)
	RawEncodeUint32(b, value)
}

func RawEncodeUint32(b *Buffer, value uint32) {
	offset := b.offset
	b.b[offset] = Uint32RawKind
	offset++
	if value < continuation {
		b.b[offset] = byte(value)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(value&(continuation-1)) | continuation
		offset++
		value >>= 7
		if value < continuation {
			b.b[offset] = byte(value)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(value&(continuation-1)) | continuation
			offset++
			value >>= 7
			if value < continuation {
				b.b[offset] = byte(value)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(value&(continuation-1)) | continuation
				offset++
				value >>= 7
				if value < continuation {
					b.b[offset] = byte(value)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(value&(continuation-1)) | continuation
					offset++
					value >>= 7
					if value < continuation {
						b.b[offset] = byte(value)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(value&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(value >> 7)
						b.offset = offset + 1
					}
				}
			}
		}
	}
}

func EncodeUint64(b *Buffer, value uint64) {
	b.grow(uint64Size)
	RawEncodeUint64(b, value)
}

func RawEncodeUint64(b *Buffer, value uint64) {
	offset := b.offset
	b.b[offset] = Uint64RawKind
	offset++
	if value < continuation {
		b.b[offset] = byte(value)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(value&(continuation-1)) | continuation
		offset++
		value >>= 7
		if value < continuation {
			b.b[offset] = byte(value)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(value&(continuation-1)) | continuation
			offset++
			value >>= 7
			if value < continuation {
				b.b[offset] = byte(value)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(value&(continuation-1)) | continuation
				offset++
				value >>= 7
				if value < continuation {
					b.b[offset] = byte(value)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(value&(continuation-1)) | continuation
					offset++
					value >>= 7
					if value < continuation {
						b.b[offset] = byte(value)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(value&(continuation-1)) | continuation
						offset++
						value >>= 7
						if value < continuation {
							b.b[offset] = byte(value)
							b.offset = offset + 1
						} else {
							b.b[offset] = byte(value&(continuation-1)) | continuation
							offset++
							value >>= 7
							if value < continuation {
								b.b[offset] = byte(value)
								b.offset = offset + 1
							} else {
								b.b[offset] = byte(value&(continuation-1)) | continuation
								offset++
								value >>= 7
								if value < continuation {
									b.b[offset] = byte(value)
									b.offset = offset + 1
								} else {
									b.b[offset] = byte(value&(continuation-1)) | continuation
									offset++
									value >>= 7
									if value < continuation {
										b.b[offset] = byte(value)
										b.offset = offset + 1
									} else {
										b.b[offset] = byte(value&(continuation-1)) | continuation
										offset++
										b.b[offset] = byte(value >> 7)
										b.offset = offset + 1
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func EncodeInt32(b *Buffer, value int32) {
	b.grow(uint32Size)
	RawEncodeInt32(b, value)
}

func RawEncodeInt32(b *Buffer, value int32) {
	castValue := uint32(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	offset := b.offset
	b.b[offset] = Int32RawKind
	offset++
	if castValue < continuation {
		b.b[offset] = byte(castValue)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(castValue&(continuation-1)) | continuation
		offset++
		castValue >>= 7
		if castValue < continuation {
			b.b[offset] = byte(castValue)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(castValue&(continuation-1)) | continuation
			offset++
			castValue >>= 7
			if castValue < continuation {
				b.b[offset] = byte(castValue)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(castValue&(continuation-1)) | continuation
				offset++
				castValue >>= 7
				if castValue < continuation {
					b.b[offset] = byte(castValue)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(castValue&(continuation-1)) | continuation
					offset++
					castValue >>= 7
					if castValue < continuation {
						b.b[offset] = byte(castValue)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(castValue&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(castValue >> 7)
						b.offset = offset + 1
					}
				}
			}
		}
	}
}

func EncodeInt64(b *Buffer, value int64) {
	b.grow(uint64Size)
	RawEncodeInt64(b, value)
}

func RawEncodeInt64(b *Buffer, value int64) {
	castValue := uint64(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	offset := b.offset
	b.b[offset] = Int64RawKind
	offset++
	if castValue < continuation {
		b.b[offset] = byte(castValue)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(castValue&(continuation-1)) | continuation
		offset++
		castValue >>= 7
		if castValue < continuation {
			b.b[offset] = byte(castValue)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(castValue&(continuation-1)) | continuation
			offset++
			castValue >>= 7
			if castValue < continuation {
				b.b[offset] = byte(castValue)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(castValue&(continuation-1)) | continuation
				offset++
				castValue >>= 7
				if castValue < continuation {
					b.b[offset] = byte(castValue)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(castValue&(continuation-1)) | continuation
					offset++
					castValue >>= 7
					if castValue < continuation {
						b.b[offset] = byte(castValue)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(castValue&(continuation-1)) | continuation
						offset++
						castValue >>= 7
						if castValue < continuation {
							b.b[offset] = byte(castValue)
							b.offset = offset + 1
						} else {
							b.b[offset] = byte(castValue&(continuation-1)) | continuation
							offset++
							castValue >>= 7
							if castValue < continuation {
								b.b[offset] = byte(castValue)
								b.offset = offset + 1
							} else {
								b.b[offset] = byte(castValue&(continuation-1)) | continuation
								offset++
								castValue >>= 7
								if castValue < continuation {
									b.b[offset] = byte(castValue)
									b.offset = offset + 1
								} else {
									b.b[offset] = byte(castValue&(continuation-1)) | continuation
									offset++
									castValue >>= 7
									if castValue < continuation {
										b.b[offset] = byte(castValue)
										b.offset = offset + 1
									} else {
										b.b[offset] = byte(castValue&(continuation-1)) | continuation
										offset++
										b.b[offset] = byte(castValue >> 7)
										b.offset = offset + 1
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func EncodeFloat32(b *Buffer, value float32) {
	b.grow(float32Size)
	RawEncodeFloat32(b, value)
}

func RawEncodeFloat32(b *Buffer, value float32) {
	offset := b.offset
	b.b[offset] = Float32RawKind
	offset++
	castValue := math.Float32bits(value)
	b.b[offset] = byte(castValue >> 24)
	offset++
	b.b[offset] = byte(castValue >> 16)
	offset++
	b.b[offset] = byte(castValue >> 8)
	offset++
	b.b[offset] = byte(castValue)
	b.offset = offset + 1
}

func EncodeFloat64(b *Buffer, value float64) {
	b.grow(float64Size)
	RawEncodeFloat64(b, value)
}

func RawEncodeFloat64(b *Buffer, value float64) {
	offset := b.offset
	b.b[offset] = Float64RawKind
	offset++
	castValue := math.Float64bits(value)
	b.b[offset] = byte(castValue >> 56)
	offset++
	b.b[offset] = byte(castValue >> 48)
	offset++
	b.b[offset] = byte(castValue >> 40)
	offset++
	b.b[offset] = byte(castValue >> 32)
	offset++
	b.b[offset] = byte(castValue >> 24)
	offset++
	b.b[offset] = byte(castValue >> 16)
	offset++
	b.b[offset] = byte(castValue >> 8)
	offset++
	b.b[offset] = byte(castValue)
	b.offset = offset + 1
}
