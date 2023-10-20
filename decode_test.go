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
	"github.com/stretchr/testify/assert"

	"errors"
	"testing"
)

func TestDecodeNil(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeNil(p)

	var value bool

	remaining, value := decodeNil(p.Bytes())
	assert.True(t, value)
	assert.Equal(t, 0, len(remaining))

	_, value = decodeNil((p.Bytes())[1:])
	assert.False(t, value)

	remaining, value = decodeNil(p.Bytes())
	assert.True(t, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.True(t, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeNil(p)
		_, _ = decodeNil(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeMap(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeMap(p, 32, StringKind, Uint32Kind)

	remaining, size, err := decodeMap(p.Bytes(), StringKind, Uint32Kind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(32), size)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeMap((p.Bytes())[1:], StringKind, Uint32Kind)
	assert.ErrorIs(t, err, InvalidMap)

	_, _, err = decodeMap(p.Bytes(), StringKind, Float64Kind)
	assert.ErrorIs(t, err, InvalidMap)

	remaining, size, err = decodeMap(p.Bytes(), StringKind, Uint32Kind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(32), size)
	assert.Equal(t, 0, len(remaining))

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeNil(p)
		remaining, size, err = decodeMap(p.Bytes(), StringKind, Uint32Kind)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeBytes(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := []byte("Test Bytes")
	encodeBytes(p, v)

	var value []byte

	remaining, value, err := decodeBytes(p.Bytes(), value)
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, value, err = decodeBytes((p.Bytes())[1:], value)
	assert.ErrorIs(t, err, InvalidBytes)

	remaining, value, err = decodeBytes(p.Bytes(), value)
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeBytes(p, v)
		remaining, value, err = decodeBytes(p.Bytes(), value)
		p.Reset()
	})
	assert.Zero(t, n)

	n = testing.AllocsPerRun(100, func() {
		encodeBytes(p, v)
		remaining, value, err = decodeBytes(p.Bytes(), nil)
		p.Reset()
	})
	assert.Equal(t, float64(1), n)

	s := [][]byte{v, v, v, v, v}
	encodeSlice(p, uint32(len(s)), BytesKind)
	for _, sb := range s {
		encodeBytes(p, sb)
	}
	var size uint32

	remaining, size, err = decodeSlice(p.Bytes(), BytesKind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(len(s)), size)

	sValue := make([][]byte, size)
	for i := uint32(0); i < size; i++ {
		remaining, sValue[i], err = decodeBytes(remaining, nil)
		assert.NoError(t, err)
		assert.Equal(t, s[i], sValue[i])
	}

	assert.Equal(t, s, sValue)
	assert.Equal(t, 0, len(remaining))

}

func TestDecodeString(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := "Test String"
	encodeString(p, v)

	var value string

	remaining, value, err := decodeString(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeString((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidString)

	remaining, value, err = decodeString(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeString(p, v)
		remaining, value, err = decodeString(p.Bytes())
		p.Reset()
	})
	assert.Equal(t, float64(1), n)

	s := []string{v, v, v, v, v}
	encodeSlice(p, uint32(len(s)), StringKind)
	for _, sb := range s {
		encodeString(p, sb)
	}
	var size uint32

	remaining, size, err = decodeSlice(p.Bytes(), StringKind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(len(s)), size)

	sValue := make([]string, size)
	for i := uint32(0); i < size; i++ {
		remaining, sValue[i], err = decodeString(remaining)
		assert.NoError(t, err)
		assert.Equal(t, s[i], sValue[i])
	}

	assert.Equal(t, s, sValue)
	assert.Equal(t, 0, len(remaining))

}

func TestDecodeError(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := errors.New("test error")
	encodeError(p, v)

	var value error

	remaining, value, err := decodeError(p.Bytes())
	assert.NoError(t, err)
	assert.ErrorIs(t, value, v)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeError((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidError)

	remaining, value, err = decodeError(p.Bytes())
	assert.NoError(t, err)
	assert.ErrorIs(t, value, v)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.ErrorIs(t, value, v)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeError(p, v)
		remaining, value, err = decodeError(p.Bytes())
		p.Reset()
	})
	assert.Equal(t, float64(2), n)

	s := []error{v, v, v, v, v}
	encodeSlice(p, uint32(len(s)), ErrorKind)
	for _, sb := range s {
		encodeError(p, sb)
	}
	var size uint32

	remaining, size, err = decodeSlice(p.Bytes(), ErrorKind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(len(s)), size)

	sValue := make([]error, size)
	for i := uint32(0); i < size; i++ {
		remaining, sValue[i], err = decodeError(remaining)
		assert.NoError(t, err)
		assert.ErrorIs(t, sValue[i], s[i])
	}

	assert.Equal(t, 0, len(remaining))

}

func TestDecodeBool(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeBool(p, true)

	var value bool

	remaining, value, err := decodeBool(p.Bytes())
	assert.NoError(t, err)
	assert.True(t, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeBool((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidBool)

	remaining, value, err = decodeBool(p.Bytes())
	assert.NoError(t, err)
	assert.True(t, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.True(t, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeBool(p, true)
		remaining, value, err = decodeBool(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)

	s := []bool{true, true, false, true, true}
	encodeSlice(p, uint32(len(s)), BoolKind)
	for _, sb := range s {
		encodeBool(p, sb)
	}
	var size uint32

	remaining, size, err = decodeSlice(p.Bytes(), BoolKind)
	assert.NoError(t, err)
	assert.Equal(t, uint32(len(s)), size)

	sValue := make([]bool, size)
	for i := uint32(0); i < size; i++ {
		remaining, sValue[i], err = decodeBool(remaining)
		assert.NoError(t, err)
		assert.Equal(t, s[i], sValue[i])
	}

	assert.Equal(t, s, sValue)
	assert.Equal(t, 0, len(remaining))

}

func TestDecodeUint8(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint8(32)
	encodeUint8(p, v)

	var value uint8

	remaining, value, err := decodeUint8(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeUint8((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidUint8)

	remaining, value, err = decodeUint8(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint8(p, v)
		remaining, value, err = decodeUint8(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeUint16(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint16(1024)
	encodeUint16(p, v)

	var value uint16

	remaining, value, err := decodeUint16(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeUint16((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidUint16)

	remaining, value, err = decodeUint16(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint16(p, v)
		remaining, value, err = decodeUint16(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeUint32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint32(4294967290)
	encodeUint32(p, v)

	var value uint32

	remaining, value, err := decodeUint32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeUint32((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidUint32)

	remaining, value, err = decodeUint32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint32(p, v)
		remaining, value, err = decodeUint32(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)

}

func TestDecodeUint64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint64(18446744073709551610)
	encodeUint64(p, v)

	var value uint64

	remaining, value, err := decodeUint64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeUint64((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidUint64)

	remaining, value, err = decodeUint64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint64(p, v)
		remaining, value, err = decodeUint64(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeInt32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := int32(2147483647)
	encodeInt32(p, v)

	var value int32

	remaining, value, err := decodeInt32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	v = int32(-2147483647)
	p.Reset()
	encodeInt32(p, v)

	remaining, value, err = decodeInt32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeInt32((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidInt32)

	remaining, value, err = decodeInt32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeInt32(p, v)
		remaining, value, err = decodeInt32(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)

}

func TestDecodeInt64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := int64(9223372036854775807)
	encodeInt64(p, v)

	var value int64

	remaining, value, err := decodeInt64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	v = int64(-9223372036854775807)
	p.Reset()
	encodeInt64(p, v)

	remaining, value, err = decodeInt64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeInt64((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidInt64)

	remaining, value, err = decodeInt64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeInt64(p, v)
		remaining, value, err = decodeInt64(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)

}

func TestDecodeFloat32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := float32(-12311.12429)
	encodeFloat32(p, v)

	var value float32

	remaining, value, err := decodeFloat32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeFloat32((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidFloat32)

	remaining, value, err = decodeFloat32(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeFloat32(p, v)
		remaining, value, err = decodeFloat32(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestDecodeFloat64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := -12311241.1242009
	encodeFloat64(p, v)

	var value float64

	remaining, value, err := decodeFloat64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	_, _, err = decodeFloat64((p.Bytes())[1:])
	assert.ErrorIs(t, err, InvalidFloat64)

	remaining, value, err = decodeFloat64(p.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, v, value)
	assert.Equal(t, 0, len(remaining))

	(p.Bytes())[len(p.Bytes())-1] = 'S'
	assert.Equal(t, v, value)

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeFloat64(p, v)
		remaining, value, err = decodeFloat64(p.Bytes())
		p.Reset()
	})
	assert.Zero(t, n)
}
