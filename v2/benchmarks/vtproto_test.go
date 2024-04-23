//go:build vtproto

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
	vtBenchmark "benchmark/vtproto/benchmark"
	"crypto/rand"
	"math"
	"runtime"
	"testing"
)

func BenchmarkEncodeVTProto(b *testing.B) {
	b.Run("Uint32", func(b *testing.B) {
		vtData := vtBenchmark.U32Data{
			U32: math.MaxUint32,
		}
		vtBuf := make([]byte, 0, 1024)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Int32", func(b *testing.B) {
		vtData := vtBenchmark.I32Data{
			I32: math.MaxInt32,
		}
		vtBuf := make([]byte, 0, 1024)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Uint64", func(b *testing.B) {
		vtData := vtBenchmark.U64Data{
			U64: math.MaxUint64,
		}
		vtBuf := make([]byte, 0, 1024)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Int64", func(b *testing.B) {
		vtData := vtBenchmark.I64Data{
			I64: math.MaxInt64,
		}
		vtBuf := make([]byte, 0, 1024)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Bytes", func(b *testing.B) {
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		vtData := vtBenchmark.BytesData{
			Bytes: randData,
		}
		vtBuf := make([]byte, 0, 1024)
		b.SetBytes(512)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Bytes (Parallel)", func(b *testing.B) {
		if testing.Short() {
			b.Skip("skipping in short mode")
		}
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		vtData := vtBenchmark.BytesData{
			Bytes: randData,
		}

		b.SetBytes(512)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var err error
			vtBuf := make([]byte, 0, 1024)
			for pb.Next() {
				_, err = vtData.MarshalToVT(vtBuf)
				if err != nil {
					b.Fatal(err)
				}
				vtBuf = vtBuf[:0]
			}
			runtime.KeepAlive(vtData)
		})
	})
}

func BenchmarkDecodeVTProto(b *testing.B) {
	b.Run("Uint32", func(b *testing.B) {
		vtData := vtBenchmark.U32Data{
			U32: math.MaxUint32,
		}
		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			if vtData.U32 != math.MaxUint32 {
				b.Fail()
			}
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Int32", func(b *testing.B) {
		vtData := vtBenchmark.I32Data{
			I32: math.MaxInt32,
		}
		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			if vtData.I32 != math.MaxInt32 {
				b.Fail()
			}
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Uint64", func(b *testing.B) {
		vtData := vtBenchmark.U64Data{
			U64: math.MaxUint64,
		}
		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			if vtData.U64 != math.MaxUint64 {
				b.Fail()
			}
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Int64", func(b *testing.B) {
		vtData := vtBenchmark.I64Data{
			I64: math.MaxInt64,
		}
		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]
		var err error
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			if vtData.I64 != math.MaxInt64 {
				b.Fail()
			}
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Bytes", func(b *testing.B) {
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)
		vtData := vtBenchmark.BytesData{
			Bytes: randData,
		}

		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]
		b.SetBytes(512)
		var err error

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vtData.Bytes = nil
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
		}
		runtime.KeepAlive(vtData)
	})

	b.Run("Bytes (Parallel)", func(b *testing.B) {
		if testing.Short() {
			b.Skip("skipping in short mode")
		}
		randData := make([]byte, 512)
		_, _ = rand.Read(randData)

		vtData := vtBenchmark.BytesData{
			Bytes: randData,
		}

		vtBuf := make([]byte, 1024)
		_, _ = vtData.MarshalToVT(vtBuf)
		vtBuf = vtBuf[:vtData.SizeVT()]

		b.SetBytes(512)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var err error
			vtData := new(vtBenchmark.BytesData)
			for pb.Next() {
				vtData.Bytes = nil
				err = vtData.UnmarshalVT(vtBuf)
				if err != nil {
					b.Fatal(err)
				}
			}
			runtime.KeepAlive(vtData)
		})
	})
}
