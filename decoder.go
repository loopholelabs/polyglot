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
	"sync"
)

var decoderPool sync.Pool

type BufferDecoder []byte

func Decoder(b []byte) *BufferDecoder {
	c := (BufferDecoder)(b)
	return &c
}

func (d *BufferDecoder) Nil() (value bool) {
	*d, value = decodeNil(*d)
	return
}

func (d *BufferDecoder) Map(keyKind, valueKind Kind) (size uint32, err error) {
	*d, size, err = decodeMap(*d, keyKind, valueKind)
	return
}

func (d *BufferDecoder) Slice(kind Kind) (size uint32, err error) {
	*d, size, err = decodeSlice(*d, kind)
	return
}

func (d *BufferDecoder) Bytes(b []byte) (value []byte, err error) {
	*d, value, err = decodeBytes(*d, b)
	return
}

func (d *BufferDecoder) String() (value string, err error) {
	*d, value, err = decodeString(*d)
	return
}

func (d *BufferDecoder) Error() (value, err error) {
	*d, value, err = decodeError(*d)
	return
}

func (d *BufferDecoder) Bool() (value bool, err error) {
	*d, value, err = decodeBool(*d)
	return
}

func (d *BufferDecoder) Uint8() (value uint8, err error) {
	*d, value, err = decodeUint8(*d)
	return
}

func (d *BufferDecoder) Uint16() (value uint16, err error) {
	*d, value, err = decodeUint16(*d)
	return
}

func (d *BufferDecoder) Uint32() (value uint32, err error) {
	*d, value, err = decodeUint32(*d)
	return
}

func (d *BufferDecoder) Uint64() (value uint64, err error) {
	*d, value, err = decodeUint64(*d)
	return
}

func (d *BufferDecoder) Int32() (value int32, err error) {
	*d, value, err = decodeInt32(*d)
	return
}

func (d *BufferDecoder) Int64() (value int64, err error) {
	*d, value, err = decodeInt64(*d)
	return
}

func (d *BufferDecoder) Float32() (value float32, err error) {
	*d, value, err = decodeFloat32(*d)
	return
}

func (d *BufferDecoder) Float64() (value float64, err error) {
	*d, value, err = decodeFloat64(*d)
	return
}
