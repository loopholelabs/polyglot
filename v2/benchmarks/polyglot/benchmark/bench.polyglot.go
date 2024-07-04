// Code generated by polyglot v1.1.4, DO NOT EDIT.
// source: bench.proto

package benchmark

import (
	"errors"
	"github.com/loopholelabs/polyglot/v2"
)

var (
	NilDecode = errors.New("cannot decode into a nil root struct")
)

type BytesData struct {
	Bytes []byte
}

func NewBytesData() *BytesData {
	return &BytesData{}
}

func (x *BytesData) Error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *BytesData) Encode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Bytes(x.Bytes)
	}
}

func (x *BytesData) Decode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	return x.decode(polyglot.Decoder(b))
}

func (x *BytesData) decode(d *polyglot.BufferDecoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.Bytes, err = d.Bytes(x.Bytes)
	if err != nil {
		return err
	}
	return nil
}

type I32Data struct {
	I32 int32
}

func NewI32Data() *I32Data {
	return &I32Data{}
}

func (x *I32Data) Error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *I32Data) Encode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Int32(x.I32)
	}
}

func (x *I32Data) Decode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	return x.decode(polyglot.Decoder(b))
}

func (x *I32Data) decode(d *polyglot.BufferDecoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.I32, err = d.Int32()
	if err != nil {
		return err
	}
	return nil
}

type U32Data struct {
	U32 uint32
}

func NewU32Data() *U32Data {
	return &U32Data{}
}

func (x *U32Data) Error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *U32Data) Encode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Uint32(x.U32)
	}
}

func (x *U32Data) Decode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	return x.decode(polyglot.Decoder(b))
}

func (x *U32Data) decode(d *polyglot.BufferDecoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.U32, err = d.Uint32()
	if err != nil {
		return err
	}
	return nil
}

type I64Data struct {
	I64 int64
}

func NewI64Data() *I64Data {
	return &I64Data{}
}

func (x *I64Data) Error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *I64Data) Encode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Int64(x.I64)
	}
}

func (x *I64Data) Decode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	return x.decode(polyglot.Decoder(b))
}

func (x *I64Data) decode(d *polyglot.BufferDecoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.I64, err = d.Int64()
	if err != nil {
		return err
	}
	return nil
}

type U64Data struct {
	U64 uint64
}

func NewU64Data() *U64Data {
	return &U64Data{}
}

func (x *U64Data) Error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *U64Data) Encode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Uint64(x.U64)
	}
}

func (x *U64Data) Decode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	return x.decode(polyglot.Decoder(b))
}

func (x *U64Data) decode(d *polyglot.BufferDecoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.U64, err = d.Uint64()
	if err != nil {
		return err
	}
	return nil
}