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
	"errors"
	"math"
)

const (
	emptyString  = ""
	VarIntLen16  = 3
	VarIntLen32  = 5
	VarIntLen64  = 10
	continuation = 0x80
)

var (
	InvalidSlice   = errors.New("invalid slice encoding")
	InvalidMap     = errors.New("invalid map encoding")
	InvalidBytes   = errors.New("invalid bytes encoding")
	InvalidString  = errors.New("invalid string encoding")
	InvalidError   = errors.New("invalid error encoding")
	InvalidBool    = errors.New("invalid bool encoding")
	InvalidUint8   = errors.New("invalid uint8 encoding")
	InvalidUint16  = errors.New("invalid uint16 encoding")
	InvalidUint32  = errors.New("invalid uint32 encoding")
	InvalidUint64  = errors.New("invalid uint64 encoding")
	InvalidInt32   = errors.New("invalid int32 encoding")
	InvalidInt64   = errors.New("invalid int64 encoding")
	InvalidFloat32 = errors.New("invalid float32 encoding")
	InvalidFloat64 = errors.New("invalid float64 encoding")
)

func decodeNil(b []byte) ([]byte, bool) {
	if len(b) > 0 {
		if b[0] == NilRawKind {
			return b[1:], true
		}
	}
	return b, false
}

func decodeMap(b []byte, keyKind, valueKind Kind) ([]byte, uint32, error) {
	if len(b) > 2 {
		if b[0] == MapRawKind && b[1] == byte(keyKind) && b[2] == byte(valueKind) {
			var size uint32
			var err error
			b, size, err = decodeUint32(b[3:])
			if err != nil {
				return b, 0, InvalidMap
			}
			return b, size, nil
		}
	}
	return b, 0, InvalidMap
}

func decodeSlice(b []byte, kind Kind) ([]byte, uint32, error) {
	if len(b) > 1 {
		if b[0] == SliceRawKind && b[1] == byte(kind) {
			var size uint32
			var err error
			b, size, err = decodeUint32(b[2:])
			if err != nil {
				return b, 0, InvalidSlice
			}
			return b, size, nil
		}
	}
	return b, 0, InvalidSlice
}

func decodeBytes(b []byte, ret []byte) ([]byte, []byte, error) {
	if len(b) > 3 && b[0] == BytesRawKind && b[1] == Uint32RawKind {
		var size int
		var offset int
		cb := uint32(b[2])
		if cb < continuation {
			size = int(cb)
			offset = 3
		} else {
			x := cb & (continuation - 1)
			cb = uint32(b[3])
			if cb < continuation {
				size = int(x | (cb << 7))
				offset = 4
			} else {
				x |= (cb & (continuation - 1)) << 7
				cb = uint32(b[4])
				if cb < continuation {
					size = int(x | (cb << 14))
					offset = 5
				} else {
					x |= (cb & (continuation - 1)) << 14
					cb = uint32(b[5])
					if cb < continuation {
						size = int(x | (cb << 21))
						offset = 6
					} else {
						x |= (cb & (continuation - 1)) << 21
						cb = uint32(b[6])
						if cb < continuation {
							size = int(x | (cb << 28))
							offset = 7
						}
					}
				}
			}
		}
		if len(b)-offset > size-1 {
			return b[size+offset:], append(ret[:0], b[offset:size+offset]...), nil
		}
	}
	return b, nil, InvalidBytes
}

func decodeString(b []byte) ([]byte, string, error) {
	if len(b) > 0 {
		if b[0] == StringRawKind {
			var size uint32
			var err error
			b, size, err = decodeUint32(b[1:])
			if err != nil {
				return b, emptyString, InvalidString
			}
			if len(b) > int(size)-1 {
				return b[size:], string(b[:size]), nil
			}
		}
	}
	return b, emptyString, InvalidString
}

func decodeError(b []byte) ([]byte, error, error) {
	if len(b) > 0 {
		if b[0] == ErrorRawKind {
			var val string
			var err error
			b, val, err = decodeString(b[1:])
			if err != nil {
				return b, nil, InvalidError
			}
			return b, Error(val), nil
		}
	}
	return b, nil, InvalidError
}

func decodeBool(b []byte) ([]byte, bool, error) {
	if len(b) > 1 {
		if b[0] == BoolRawKind {
			if b[1] == trueBool {
				return b[2:], true, nil
			} else {
				return b[2:], false, nil
			}
		}
	}
	return b, false, InvalidBool
}

func decodeUint8(b []byte) ([]byte, uint8, error) {
	if len(b) > 1 {
		if b[0] == Uint8RawKind {
			return b[2:], b[1], nil
		}
	}
	return b, 0, InvalidUint8
}

func decodeUint16(b []byte) ([]byte, uint16, error) {
	if len(b) > 1 && b[0] == Uint16RawKind {
		cb := uint16(b[1])
		if cb < continuation {
			return b[2:], cb, nil
		}

		x := cb & (continuation - 1)
		cb = uint16(b[2])
		if cb < continuation {
			return b[3:], x | (cb << 7), nil
		}

		x |= (cb & (continuation - 1)) << 7
		cb = uint16(b[3])
		if cb < continuation {
			return b[4:], x | (cb << 14), nil
		}
	}
	return b, 0, InvalidUint16
}

func decodeUint32(b []byte) ([]byte, uint32, error) {
	if len(b) > 1 && b[0] == Uint32RawKind {
		cb := uint32(b[1])
		if cb < continuation {
			return b[2:], cb, nil
		}

		x := cb & (continuation - 1)
		cb = uint32(b[2])
		if cb < continuation {
			return b[3:], x | (cb << 7), nil
		}

		x |= (cb & (continuation - 1)) << 7
		cb = uint32(b[3])
		if cb < continuation {
			return b[4:], x | (cb << 14), nil
		}

		x |= (cb & (continuation - 1)) << 14
		cb = uint32(b[4])
		if cb < continuation {
			return b[5:], x | (cb << 21), nil
		}

		x |= (cb & (continuation - 1)) << 21
		cb = uint32(b[5])
		if cb < continuation {
			return b[6:], x | (cb << 28), nil
		}
	}
	return b, 0, InvalidUint32
}

func decodeUint64(b []byte) ([]byte, uint64, error) {
	if len(b) > 1 && b[0] == Uint64RawKind {
		cb := uint64(b[1])
		if cb < continuation {
			return b[2:], cb, nil
		}

		x := cb & (continuation - 1)
		cb = uint64(b[2])
		if cb < continuation {
			return b[3:], x | (cb << 7), nil
		}

		x |= (cb & (continuation - 1)) << 7
		cb = uint64(b[3])
		if cb < continuation {
			return b[4:], x | (cb << 14), nil
		}

		x |= (cb & (continuation - 1)) << 14
		cb = uint64(b[4])
		if cb < continuation {
			return b[5:], x | (cb << 21), nil
		}

		x |= (cb & (continuation - 1)) << 21
		cb = uint64(b[5])
		if cb < continuation {
			return b[6:], x | (cb << 28), nil
		}

		x |= (cb & (continuation - 1)) << 28
		cb = uint64(b[6])
		if cb < continuation {
			return b[7:], x | (cb << 35), nil
		}

		x |= (cb & (continuation - 1)) << 35
		cb = uint64(b[7])
		if cb < continuation {
			return b[8:], x | (cb << 42), nil
		}

		x |= (cb & (continuation - 1)) << 42
		cb = uint64(b[8])
		if cb < continuation {
			return b[9:], x | (cb << 49), nil
		}

		x |= (cb & (continuation - 1)) << 49
		cb = uint64(b[9])
		if cb < continuation {
			return b[10:], x | (cb << 56), nil
		}

		x |= (cb & (continuation - 1)) << 56
		cb = uint64(b[10])
		if cb < continuation {
			return b[11:], x | (cb << 63), nil
		}
	}
	return b, 0, InvalidUint64
}

func decodeInt32(b []byte) ([]byte, int32, error) {
	if len(b) > 1 && b[0] == Int32RawKind {
		cb := uint32(b[1])
		if cb < continuation {
			x := int32(cb >> 1)
			if cb&1 != 0 {
				x = -(x + 1)
			}
			return b[2:], x, nil
		}

		x := cb & (continuation - 1)
		cb = uint32(b[2])
		if cb < continuation {
			x |= cb << 7
			if x&1 != 0 {
				return b[3:], -(int32(x>>1) + 1), nil
			}
			return b[3:], int32(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 7
		cb = uint32(b[3])
		if cb < continuation {
			x |= cb << 14
			if x&1 != 0 {
				return b[4:], -(int32(x>>1) + 1), nil
			}
			return b[4:], int32(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 14
		cb = uint32(b[4])
		if cb < continuation {
			x |= cb << 21
			if x&1 != 0 {
				return b[5:], -(int32(x>>1) + 1), nil
			}
			return b[5:], int32(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 21
		cb = uint32(b[5])
		if cb < continuation {
			x |= cb << 28
			if x&1 != 0 {
				return b[6:], -(int32(x>>1) + 1), nil
			}
			return b[6:], int32(x >> 1), nil
		}
	}
	return b, 0, InvalidInt32
}

func decodeInt64(b []byte) ([]byte, int64, error) {
	if len(b) > 1 && b[0] == Int64RawKind {
		cb := uint64(b[1])
		if cb < continuation {
			x := int64(cb >> 1)
			if cb&1 != 0 {
				x = -(x + 1)
			}
			return b[2:], x, nil
		}

		x := cb & (continuation - 1)
		cb = uint64(b[2])
		if cb < continuation {
			x |= cb << 7
			if x&1 != 0 {
				return b[3:], -(int64(x>>1) + 1), nil
			}
			return b[3:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 7
		cb = uint64(b[3])
		if cb < continuation {
			x |= cb << 14
			if x&1 != 0 {
				return b[4:], -(int64(x>>1) + 1), nil
			}
			return b[4:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 14
		cb = uint64(b[4])
		if cb < continuation {
			x |= cb << 21
			if x&1 != 0 {
				return b[5:], -(int64(x>>1) + 1), nil
			}
			return b[5:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 21
		cb = uint64(b[5])
		if cb < continuation {
			x |= cb << 28
			if x&1 != 0 {
				return b[6:], -(int64(x>>1) + 1), nil
			}
			return b[6:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 28
		cb = uint64(b[6])
		if cb < continuation {
			x |= cb << 35
			if x&1 != 0 {
				return b[7:], -(int64(x>>1) + 1), nil
			}
			return b[7:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 35
		cb = uint64(b[7])
		if cb < continuation {
			x |= cb << 42
			if x&1 != 0 {
				return b[8:], -(int64(x>>1) + 1), nil
			}
			return b[8:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 42
		cb = uint64(b[8])
		if cb < continuation {
			x |= cb << 49
			if x&1 != 0 {
				return b[9:], -(int64(x>>1) + 1), nil
			}
			return b[9:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 49
		cb = uint64(b[9])
		if cb < continuation {
			x |= cb << 56
			if x&1 != 0 {
				return b[10:], -(int64(x>>1) + 1), nil
			}
			return b[10:], int64(x >> 1), nil
		}

		x |= (cb & (continuation - 1)) << 56
		cb = uint64(b[10])
		if cb < continuation {
			x |= cb << 63
			if x&1 != 0 {
				return b[11:], -(int64(x>>1) + 1), nil
			}
			return b[11:], int64(x >> 1), nil
		}
	}
	return b, 0, InvalidInt64
}

func decodeFloat32(b []byte) ([]byte, float32, error) {
	if len(b) > 4 && b[0] == Float32RawKind {
		return b[5:], math.Float32frombits(uint32(b[4]) | uint32(b[3])<<8 | uint32(b[2])<<16 | uint32(b[1])<<24), nil
	}
	return b, 0, InvalidFloat32
}

func decodeFloat64(b []byte) ([]byte, float64, error) {
	if len(b) > 8 && b[0] == Float64RawKind {
		return b[9:], math.Float64frombits(uint64(b[8]) | uint64(b[7])<<8 | uint64(b[6])<<16 | uint64(b[5])<<24 |
			uint64(b[4])<<32 | uint64(b[3])<<40 | uint64(b[2])<<48 | uint64(b[1])<<56), nil
	}
	return b, 0, InvalidFloat64
}
