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

func encodeNil(b *Buffer) {
	b.Grow(nilSize)
	b.b[b.offset] = NilRawKind
	b.offset++
}

func encodeMap(b *Buffer, size uint32, keyKind Kind, valueKind Kind) {
	b.Grow(mapSize)
	offset := b.offset
	b.b[offset] = MapRawKind
	offset++
	b.b[offset] = byte(keyKind)
	offset++
	b.b[offset] = byte(valueKind)
	offset++
	b.b[offset] = Uint32RawKind
	offset++
	if size < continuation {
		b.b[offset] = byte(size)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(size&(continuation-1)) | continuation
		offset++
		size >>= 7
		if size < continuation {
			b.b[offset] = byte(size)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(size&(continuation-1)) | continuation
			offset++
			size >>= 7
			if size < continuation {
				b.b[offset] = byte(size)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(size&(continuation-1)) | continuation
				offset++
				size >>= 7
				if size < continuation {
					b.b[offset] = byte(size)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(size&(continuation-1)) | continuation
					offset++
					size >>= 7
					if size < continuation {
						b.b[offset] = byte(size)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(size&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(size >> 7)
						b.offset = offset + 1
					}
				}
			}
		}
	}
}

func encodeSlice(b *Buffer, size uint32, kind Kind) {
	b.Grow(sliceSize)
	offset := b.offset
	b.b[offset] = SliceRawKind
	offset++
	b.b[offset] = byte(kind)
	offset++
	b.b[offset] = Uint32RawKind
	offset++
	if size < continuation {
		b.b[offset] = byte(size)
		b.offset = offset + 1
	} else {
		b.b[offset] = byte(size&(continuation-1)) | continuation
		offset++
		size >>= 7
		if size < continuation {
			b.b[offset] = byte(size)
			b.offset = offset + 1
		} else {
			b.b[offset] = byte(size&(continuation-1)) | continuation
			offset++
			size >>= 7
			if size < continuation {
				b.b[offset] = byte(size)
				b.offset = offset + 1
			} else {
				b.b[offset] = byte(size&(continuation-1)) | continuation
				offset++
				size >>= 7
				if size < continuation {
					b.b[offset] = byte(size)
					b.offset = offset + 1
				} else {
					b.b[offset] = byte(size&(continuation-1)) | continuation
					offset++
					size >>= 7
					if size < continuation {
						b.b[offset] = byte(size)
						b.offset = offset + 1
					} else {
						b.b[offset] = byte(size&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(size >> 7)
						b.offset = offset + 1
					}
				}
			}
		}
	}
}

func encodeBytes(b *Buffer, value []byte) {
	b.Grow(bytesSize + len(value))
	castValue := uint32(len(value))
	offset := b.offset
	b.b[offset] = BytesRawKind
	offset++
	b.b[offset] = Uint32RawKind
	offset++
	if castValue < continuation {
		b.b[offset] = byte(castValue)
		offset++
	} else {
		b.b[offset] = byte(castValue&(continuation-1)) | continuation
		offset++
		castValue >>= 7
		if castValue < continuation {
			b.b[offset] = byte(castValue)
			offset++
		} else {
			b.b[offset] = byte(castValue&(continuation-1)) | continuation
			offset++
			castValue >>= 7
			if castValue < continuation {
				b.b[offset] = byte(castValue)
				offset++
			} else {
				b.b[offset] = byte(castValue&(continuation-1)) | continuation
				offset++
				castValue >>= 7
				if castValue < continuation {
					b.b[offset] = byte(castValue)
					offset++
				} else {
					b.b[offset] = byte(castValue&(continuation-1)) | continuation
					offset++
					castValue >>= 7
					if castValue < continuation {
						b.b[offset] = byte(castValue)
					} else {
						b.b[offset] = byte(castValue&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(castValue >> 7)
						offset++
					}
				}
			}
		}
	}
	b.offset = offset + copy(b.b[offset:], value)
}
func encodeString(b *Buffer, value string) {
	b.Grow(stringSize + len(value))
	var nb []byte
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&nb))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&value))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	offset := b.offset
	castValue := uint32(len(nb))
	b.b[offset] = StringRawKind
	offset++
	b.b[offset] = Uint32RawKind
	offset++
	if castValue < continuation {
		b.b[offset] = byte(castValue)
		offset++
	} else {
		b.b[offset] = byte(castValue&(continuation-1)) | continuation
		offset++
		castValue >>= 7
		if castValue < continuation {
			b.b[offset] = byte(castValue)
			offset++
		} else {
			b.b[offset] = byte(castValue&(continuation-1)) | continuation
			offset++
			castValue >>= 7
			if castValue < continuation {
				b.b[offset] = byte(castValue)
				offset++
			} else {
				b.b[offset] = byte(castValue&(continuation-1)) | continuation
				offset++
				castValue >>= 7
				if castValue < continuation {
					b.b[offset] = byte(castValue)
					offset++
				} else {
					b.b[offset] = byte(castValue&(continuation-1)) | continuation
					offset++
					castValue >>= 7
					if castValue < continuation {
						b.b[offset] = byte(castValue)
						offset++
					} else {
						b.b[offset] = byte(castValue&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(castValue >> 7)
						offset++
					}
				}
			}
		}
	}
	b.offset = offset + copy(b.b[offset:], nb)
}

func encodeError(b *Buffer, err error) {
	errString := err.Error()
	b.Grow(errorSize + len(errString))
	offset := b.offset
	b.b[offset] = ErrorRawKind
	offset++

	var nb []byte
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&nb))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&errString))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	castValue := uint32(len(nb))
	b.b[offset] = StringRawKind
	offset++

	b.b[offset] = Uint32RawKind
	offset++
	if castValue < continuation {
		b.b[offset] = byte(castValue)
		offset++
	} else {
		b.b[offset] = byte(castValue&(continuation-1)) | continuation
		offset++
		castValue >>= 7
		if castValue < continuation {
			b.b[offset] = byte(castValue)
			offset++
		} else {
			b.b[offset] = byte(castValue&(continuation-1)) | continuation
			offset++
			castValue >>= 7
			if castValue < continuation {
				b.b[offset] = byte(castValue)
				offset++
			} else {
				b.b[offset] = byte(castValue&(continuation-1)) | continuation
				offset++
				castValue >>= 7
				if castValue < continuation {
					b.b[offset] = byte(castValue)
					offset++
				} else {
					b.b[offset] = byte(castValue&(continuation-1)) | continuation
					offset++
					castValue >>= 7
					if castValue < continuation {
						b.b[offset] = byte(castValue)
						offset++
					} else {
						b.b[offset] = byte(castValue&(continuation-1)) | continuation
						offset++
						b.b[offset] = byte(castValue >> 7)
						offset++
					}
				}
			}
		}
	}
	b.offset = offset + copy(b.b[offset:], nb)
}

func encodeBool(b *Buffer, value bool) {
	b.Grow(boolSize)
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

func encodeUint8(b *Buffer, value uint8) {
	b.Grow(uint8Size)
	offset := b.offset
	b.b[offset] = Uint8RawKind
	offset++
	b.b[offset] = value
	b.offset = offset + 1
}

func encodeUint16(b *Buffer, value uint16) {
	b.Grow(uint16Size)
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

func encodeUint32(b *Buffer, value uint32) {
	b.Grow(uint32Size)
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

func encodeUint64(b *Buffer, value uint64) {
	b.Grow(uint64Size)
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

func encodeInt32(b *Buffer, value int32) {
	b.Grow(uint32Size)
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

func encodeInt64(b *Buffer, value int64) {
	b.Grow(uint64Size)
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

func encodeFloat32(b *Buffer, value float32) {
	b.Grow(float32Size)
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

func encodeFloat64(b *Buffer, value float64) {
	b.Grow(float64Size)
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
