//go:build !vtproto

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

package benchmarks

import (
	polyglotBenchmark "benchmark/polyglot/benchmark"
	"crypto/rand"
	"github.com/loopholelabs/polyglot"
	"math"
	"runtime"
	"testing"
)

func BenchmarkEncodePolyglot(b *testing.B) {
	b.Run("Uint32", func(b *testing.B) {
		polyglotData := polyglotBenchmark.U32Data{
			U32: math.MaxUint32,
		}
		polyglotBuf := polyglot.NewBuffer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Int32", func(b *testing.B) {
		polyglotData := polyglotBenchmark.I32Data{
			I32: math.MaxInt32,
		}
		polyglotBuf := polyglot.NewBuffer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Uint64", func(b *testing.B) {
		polyglotData := polyglotBenchmark.U64Data{
			U64: math.MaxUint64,
		}
		polyglotBuf := polyglot.NewBuffer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Int64", func(b *testing.B) {
		polyglotData := polyglotBenchmark.I64Data{
			I64: math.MaxInt64,
		}
		polyglotBuf := polyglot.NewBuffer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Bytes", func(b *testing.B) {
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		polyglotData := polyglotBenchmark.BytesData{
			Bytes: randData,
		}
		polyglotBuf := polyglot.NewBuffer()
		b.SetBytes(512)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Bytes (Parallel)", func(b *testing.B) {
		if testing.Short() {
			b.Skip("skipping in short mode")
		}

		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		polyglotData := polyglotBenchmark.BytesData{
			Bytes: randData,
		}
		b.SetBytes(512)
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			polyglotBuf := polyglot.NewBufferSize(1024)
			for pb.Next() {
				polyglotData.Encode(polyglotBuf)
				polyglotBuf.Reset()
			}
			runtime.KeepAlive(polyglotData)
		})
	})
}

func BenchmarkDecodePolyglot(b *testing.B) {
	b.Run("Uint32", func(b *testing.B) {
		polyglotData := polyglotBenchmark.U32Data{
			U32: math.MaxUint32,
		}
		polyglotBuf := polyglot.NewBuffer()
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		var err error
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
			if polyglotData.U32 != math.MaxUint32 {
				b.Fail()
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Int32", func(b *testing.B) {
		polyglotData := polyglotBenchmark.I32Data{
			I32: math.MaxInt32,
		}
		polyglotBuf := polyglot.NewBuffer()
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		var err error
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
			if polyglotData.I32 != math.MaxInt32 {
				b.Fail()
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Uint64", func(b *testing.B) {
		polyglotData := polyglotBenchmark.U64Data{
			U64: math.MaxUint64,
		}
		polyglotBuf := polyglot.NewBuffer()
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
			if polyglotData.U64 != math.MaxUint64 {
				b.Fail()
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Int64", func(b *testing.B) {
		polyglotData := polyglotBenchmark.I64Data{
			I64: math.MaxInt64,
		}
		polyglotBuf := polyglot.NewBuffer()
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		var err error
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
			if polyglotData.I64 != math.MaxInt64 {
				b.Fail()
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Bytes", func(b *testing.B) {
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		polyglotData := polyglotBenchmark.BytesData{
			Bytes: randData,
		}
		polyglotBuf := polyglot.NewBuffer()
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		var err error
		b.SetBytes(512)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			polyglotData.Bytes = nil
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("Bytes (Parallel)", func(b *testing.B) {
		if testing.Short() {
			b.Skip("skipping in short mode")
		}
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		polyglotData := polyglotBenchmark.BytesData{
			Bytes: randData,
		}
		polyglotBuf := polyglot.NewBufferSize(1024)
		polyglotData.Encode(polyglotBuf)
		polyglotBytes := polyglotBuf.Bytes()
		b.SetBytes(512)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var err error
			polyglotData := new(polyglotBenchmark.BytesData)
			for pb.Next() {
				polyglotData.Bytes = nil
				err = polyglotData.Decode(polyglotBytes)
				if err != nil {
					b.Fatal(err)
				}
			}
			runtime.KeepAlive(polyglotData)
		})
	})
}
