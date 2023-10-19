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
	b.Grow(nilSize)
	RawEncodeNil(b)
}

func RawEncodeNil(b *Buffer) {
	b.WriteRawByte(NilRawKind)
}

func EncodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	b.Grow(mapSize)
	RawEncodeMap(b, size, keyKind, valueKind)
}

func RawEncodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	b.WriteRawByteDirect(MapRawKind, 0)
	b.WriteRawByteDirect(byte(keyKind), 1)
	b.WriteRawByteDirect(byte(valueKind), 2)
	b.AddOffset(3)
	RawEncodeUint32(b, size)
}

func EncodeSlice(b *Buffer, size uint32, kind Kind) {
	b.Grow(sliceSize)
	RawEncodeSlice(b, size, kind)
}

func RawEncodeSlice(b *Buffer, size uint32, kind Kind) {
	b.WriteRawByte(SliceRawKind)
	b.WriteRawByte(byte(kind))
	RawEncodeUint32(b, size)
}

func EncodeBytes(b *Buffer, value []byte) {
	b.Grow(bytesSize + len(value))
	RawEncodeBytes(b, value)
}

func RawEncodeBytes(b *Buffer, value []byte) {
	b.WriteRawByte(BytesRawKind)
	RawEncodeUint32(b, uint32(len(value)))
	b.WriteRaw(value)
}

func EncodeString(b *Buffer, value string) {
	b.Grow(stringSize + len(value))
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
	b.WriteRawByte(StringRawKind)
	RawEncodeUint32(b, uint32(len(nb)))
	b.WriteRaw(nb)
}

func EncodeError(b *Buffer, err error) {
	errString := err.Error()
	b.Grow(errorSize + len(errString))
	RawEncodeError(b, errString)
}

func RawEncodeError(b *Buffer, errString string) {
	b.WriteRawByte(ErrorRawKind)
	RawEncodeString(b, errString)
}

func EncodeBool(b *Buffer, value bool) {
	b.Grow(boolSize)
	RawEncodeBool(b, value)
}

func RawEncodeBool(b *Buffer, value bool) {
	b.WriteRawByte(BoolRawKind)
	if value {
		b.WriteRawByte(trueBool)
	} else {
		b.WriteRawByte(falseBool)
	}
}

func EncodeUint8(b *Buffer, value uint8) {
	b.Grow(uint8Size)
	RawEncodeUint8(b, value)
}

func RawEncodeUint8(b *Buffer, value uint8) {
	b.WriteRawByte(Uint8RawKind)
	b.WriteRawByte(value)
}

func EncodeUint16(b *Buffer, value uint16) {
	b.Grow(uint16Size)
	RawEncodeUint16(b, value)
}

func RawEncodeUint16(b *Buffer, value uint16) {
	b.WriteRawByte(Uint16RawKind)
	if value < (1 << 7) {
		b.WriteRawByte(byte(value))
	} else {
		b.WriteRawByte(byte(value&(continuation-1)) | continuation)
		value >>= 7
		if value < (1 << 7) {
			b.WriteRawByte(byte(value))
		} else {
			b.WriteRawByte(byte(value&(continuation-1)) | continuation)
			value >>= 7
			if value < (1 << 7) {
				b.WriteRawByte(byte(value))
			} else {
				b.WriteRawByte(byte(value&(continuation-1)) | continuation)
				b.WriteRawByte(byte(value >> 7))
			}
		}
	}
}

func EncodeUint32(b *Buffer, value uint32) {
	b.Grow(uint32Size)
	RawEncodeUint32(b, value)
}

func RawEncodeUint32(b *Buffer, value uint32) {
	b.WriteRawByteDirect(Uint32RawKind, 0)
	if value < (1 << 7) {
		b.WriteRawByteDirect(byte(value), 1)
		b.AddOffset(2)
	} else {
		b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 1)
		value >>= 7
		if value < (1 << 7) {
			b.WriteRawByteDirect(byte(value), 2)
			b.AddOffset(3)
		} else {
			b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 2)
			value >>= 7
			if value < (1 << 7) {
				b.WriteRawByteDirect(byte(value), 3)
				b.AddOffset(4)
			} else {
				b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 3)
				value >>= 7
				if value < (1 << 7) {
					b.WriteRawByteDirect(byte(value), 4)
					b.AddOffset(5)
				} else {
					b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 4)
					value >>= 7
					if value < (1 << 7) {
						b.WriteRawByteDirect(byte(value), 5)
						b.AddOffset(6)
					} else {
						b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 5)
						b.WriteRawByteDirect(byte(value>>7), 6)
						b.AddOffset(7)
					}
				}
			}
		}
	}
}

func EncodeUint64(b *Buffer, value uint64) {
	b.Grow(uint64Size)
	RawEncodeUint64(b, value)
}

func RawEncodeUint64(b *Buffer, value uint64) {
	b.WriteRawByteDirect(Uint64RawKind, 0)
	if value < (1 << 7) {
		b.WriteRawByteDirect(byte(value), 1)
		b.AddOffset(2)
	} else {
		b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 1)
		value >>= 7
		if value < (1 << 7) {
			b.WriteRawByteDirect(byte(value), 2)
			b.AddOffset(3)
		} else {
			b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 2)
			value >>= 7
			if value < (1 << 7) {
				b.WriteRawByteDirect(byte(value), 3)
				b.AddOffset(4)
			} else {
				b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 3)
				value >>= 7
				if value < (1 << 7) {
					b.WriteRawByteDirect(byte(value), 4)
					b.AddOffset(5)
				} else {
					b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 4)
					value >>= 7
					if value < (1 << 7) {
						b.WriteRawByteDirect(byte(value), 5)
						b.AddOffset(6)
					} else {
						b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 5)
						value >>= 7
						if value < (1 << 7) {
							b.WriteRawByteDirect(byte(value), 6)
							b.AddOffset(7)
						} else {
							b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 6)
							value >>= 7
							if value < (1 << 7) {
								b.WriteRawByteDirect(byte(value), 7)
								b.AddOffset(8)
							} else {
								b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 7)
								value >>= 7
								if value < (1 << 7) {
									b.WriteRawByteDirect(byte(value), 8)
									b.AddOffset(9)
								} else {
									b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 8)
									value >>= 7
									if value < (1 << 7) {
										b.WriteRawByteDirect(byte(value), 9)
										b.AddOffset(10)
									} else {
										b.WriteRawByteDirect(byte(value&(continuation-1))|continuation, 9)
										b.WriteRawByteDirect(byte(value>>7), 10)
										b.AddOffset(11)
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
	b.Grow(uint32Size)
	RawEncodeInt32(b, value)
}

func RawEncodeInt32(b *Buffer, value int32) {
	castValue := uint32(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	b.WriteRawByteDirect(Int32RawKind, 0)
	if castValue < (1 << 7) {
		b.WriteRawByteDirect(byte(castValue), 1)
		b.AddOffset(2)
	} else {
		b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 1)
		castValue >>= 7
		if castValue < (1 << 7) {
			b.WriteRawByteDirect(byte(castValue), 2)
			b.AddOffset(3)
		} else {
			b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 2)
			castValue >>= 7
			if castValue < (1 << 7) {
				b.WriteRawByteDirect(byte(castValue), 3)
				b.AddOffset(4)
			} else {
				b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 3)
				castValue >>= 7
				if castValue < (1 << 7) {
					b.WriteRawByteDirect(byte(castValue), 4)
					b.AddOffset(5)
				} else {
					b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 4)
					castValue >>= 7
					if castValue < (1 << 7) {
						b.WriteRawByteDirect(byte(castValue), 5)
						b.AddOffset(6)
					} else {
						b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 5)
						b.WriteRawByteDirect(byte(castValue>>7), 6)
						b.AddOffset(7)
					}
				}
			}
		}
	}
}

func EncodeInt64(b *Buffer, value int64) {
	b.Grow(uint64Size)
	RawEncodeInt64(b, value)
}

func RawEncodeInt64(b *Buffer, value int64) {
	castValue := uint64(value) << 1
	if value < 0 {
		castValue = ^castValue
	}
	b.WriteRawByteDirect(Int64RawKind, 0)
	if castValue < (1 << 7) {
		b.WriteRawByteDirect(byte(castValue), 1)
		b.AddOffset(2)
	} else {
		b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 1)
		castValue >>= 7
		if castValue < (1 << 7) {
			b.WriteRawByteDirect(byte(castValue), 2)
			b.AddOffset(3)
		} else {
			b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 2)
			castValue >>= 7
			if castValue < (1 << 7) {
				b.WriteRawByteDirect(byte(castValue), 3)
				b.AddOffset(4)
			} else {
				b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 3)
				castValue >>= 7
				if castValue < (1 << 7) {
					b.WriteRawByteDirect(byte(castValue), 4)
					b.AddOffset(5)
				} else {
					b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 4)
					castValue >>= 7
					if castValue < (1 << 7) {
						b.WriteRawByteDirect(byte(castValue), 5)
						b.AddOffset(6)
					} else {
						b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 5)
						castValue >>= 7
						if castValue < (1 << 7) {
							b.WriteRawByteDirect(byte(castValue), 6)
							b.AddOffset(7)
						} else {
							b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 6)
							castValue >>= 7
							if castValue < (1 << 7) {
								b.WriteRawByteDirect(byte(castValue), 7)
								b.AddOffset(8)
							} else {
								b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 7)
								castValue >>= 7
								if castValue < (1 << 7) {
									b.WriteRawByteDirect(byte(castValue), 8)
									b.AddOffset(9)
								} else {
									b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 8)
									castValue >>= 7
									if castValue < (1 << 7) {
										b.WriteRawByteDirect(byte(castValue), 9)
										b.AddOffset(10)
									} else {
										b.WriteRawByteDirect(byte(castValue&(continuation-1))|continuation, 9)
										b.WriteRawByteDirect(byte(castValue>>7), 10)
										b.AddOffset(11)
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
	b.Grow(float32Size)
	RawEncodeFloat32(b, value)
}

func RawEncodeFloat32(b *Buffer, value float32) {
	b.WriteRawByte(Float32RawKind)
	castValue := math.Float32bits(value)
	b.WriteRawByte(byte(castValue >> 24))
	b.WriteRawByte(byte(castValue >> 16))
	b.WriteRawByte(byte(castValue >> 8))
	b.WriteRawByte(byte(castValue))
}

func EncodeFloat64(b *Buffer, value float64) {
	b.Grow(float64Size)
	RawEncodeFloat64(b, value)
}

func RawEncodeFloat64(b *Buffer, value float64) {
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
