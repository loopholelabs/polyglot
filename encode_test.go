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

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestEncodeNil(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeNil(p)

	assert.Equal(t, 1, len(p.Bytes()))
	assert.Equal(t, NilKind, Kind(p.Bytes()))

	n := testing.AllocsPerRun(100, func() {
		encodeNil(p)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeMap(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeMap(p, 32, StringKind, Uint32Kind)

	assert.Equal(t, 1+1+1+1+4, len(p.Bytes()))
	assert.Equal(t, MapKind, Kind((p.Bytes())[0:1]))
	assert.Equal(t, StringKind, Kind((p.Bytes())[1:2]))
	assert.Equal(t, Uint32Kind, Kind((p.Bytes())[2:3]))
	assert.Equal(t, Uint32Kind, Kind((p.Bytes())[3:4]))

	n := testing.AllocsPerRun(100, func() {
		encodeMap(p, 32, StringKind, Uint32Kind)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeSlice(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	encodeSlice(p, 32, StringKind)

	assert.Equal(t, 1+1+1+4, len(p.Bytes()))
	assert.Equal(t, SliceKind, Kind((p.Bytes())[0:1]))
	assert.Equal(t, StringKind, Kind((p.Bytes())[1:2]))
	assert.Equal(t, Uint32Kind, Kind((p.Bytes())[2:3]))

	n := testing.AllocsPerRun(100, func() {
		encodeSlice(p, 32, StringKind)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeBytes(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := []byte("Test String")

	encodeBytes(p, v)

	assert.Equal(t, 1+1+4+len(v), len(p.Bytes()))
	assert.Equal(t, v, (p.Bytes())[1+1+4:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeBytes(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeString(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := "Test String"
	e := []byte(v)

	encodeString(p, v)

	assert.Equal(t, 1+1+4+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1+1+4:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeString(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeError(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := errors.New("Test Error")
	e := []byte(v.Error())

	encodeError(p, v)

	assert.Equal(t, 1+1+1+4+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1+1+1+4:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeError(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeBool(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	e := []byte{trueBool}

	encodeBool(p, true)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeBool(p, true)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeUint8(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint8(32)
	e := []byte{v}

	encodeUint8(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint8(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeUint16(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint16(1024)
	e := []byte{byte(v >> 8), byte(v)}

	encodeUint16(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint16(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeUint32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint32(4294967290)
	e := []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}

	encodeUint32(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint32(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeUint64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := uint64(18446744073709551610)
	e := []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}

	encodeUint64(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeUint64(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeInt32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := int32(-2147483648)
	e := []byte{byte(uint32(v) >> 24), byte(uint32(v) >> 16), byte(uint32(v) >> 8), byte(uint32(v))}

	encodeInt32(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeInt32(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeInt64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := int64(-9223372036854775808)
	e := []byte{byte(uint64(v) >> 56), byte(uint64(v) >> 48), byte(uint64(v) >> 40), byte(uint64(v) >> 32), byte(uint64(v) >> 24), byte(uint64(v) >> 16), byte(uint64(v) >> 8), byte(uint64(v))}

	encodeInt64(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeInt64(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeFloat32(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := float32(-214648.34432)
	e := []byte{byte(math.Float32bits(v) >> 24), byte(math.Float32bits(v) >> 16), byte(math.Float32bits(v) >> 8), byte(math.Float32bits(v))}

	encodeFloat32(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeFloat32(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}

func TestEncodeFloat64(t *testing.T) {
	t.Parallel()

	p := NewBuffer()
	v := -922337203685.2345
	e := []byte{byte(math.Float64bits(v) >> 56), byte(math.Float64bits(v) >> 48), byte(math.Float64bits(v) >> 40), byte(math.Float64bits(v) >> 32), byte(math.Float64bits(v) >> 24), byte(math.Float64bits(v) >> 16), byte(math.Float64bits(v) >> 8), byte(math.Float64bits(v))}

	encodeFloat64(p, v)

	assert.Equal(t, 1+len(e), len(p.Bytes()))
	assert.Equal(t, e, (p.Bytes())[1:])

	p.Reset()
	n := testing.AllocsPerRun(100, func() {
		encodeFloat64(p, v)
		p.Reset()
	})
	assert.Zero(t, n)
}
