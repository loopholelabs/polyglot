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

package demo

import (
	polyglotBenchmark "benchmark/polyglot/benchmark"
	vtBenchmark "benchmark/vtproto/benchmark"
	"crypto/rand"
	"github.com/loopholelabs/polyglot"
	"math"
	"runtime"
	"testing"
)

func BenchmarkEncodeUint32(b *testing.B) {
	polyglotData := polyglotBenchmark.U32Data{
		U32: math.MaxUint32,
	}
	polyglotBuf := polyglot.NewBuffer()

	vtData := vtBenchmark.U32Data{
		U32: math.MaxUint32,
	}
	vtBuf := make([]byte, 0, 1024)

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkDecodeUint32(b *testing.B) {
	polyglotData := polyglotBenchmark.U32Data{
		U32: math.MaxUint32,
	}
	polyglotBuf := polyglot.NewBuffer()
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.U32Data{
		U32: math.MaxUint32,
	}
	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		var err error
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

	b.Run("vtproto", func(b *testing.B) {
		var err error
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
}

func BenchmarkEncodeInt32(b *testing.B) {
	polyglotData := polyglotBenchmark.I32Data{
		I32: math.MaxInt32,
	}
	polyglotBuf := polyglot.NewBuffer()

	vtData := vtBenchmark.I32Data{
		I32: math.MaxInt32,
	}
	vtBuf := make([]byte, 0, 1024)

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkDecodeInt32(b *testing.B) {
	polyglotData := polyglotBenchmark.I32Data{
		I32: math.MaxInt32,
	}
	polyglotBuf := polyglot.NewBuffer()
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.I32Data{
		I32: math.MaxInt32,
	}
	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		var err error
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

	b.Run("vtproto", func(b *testing.B) {
		var err error
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
}

func BenchmarkEncodeUint64(b *testing.B) {
	polyglotData := polyglotBenchmark.U64Data{
		U64: math.MaxUint64,
	}
	polyglotBuf := polyglot.NewBuffer()

	vtData := vtBenchmark.U64Data{
		U64: math.MaxUint64,
	}
	vtBuf := make([]byte, 0, 1024)

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkDecodeUint64(b *testing.B) {
	polyglotData := polyglotBenchmark.U64Data{
		U64: math.MaxUint64,
	}
	polyglotBuf := polyglot.NewBuffer()
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.U64Data{
		U64: math.MaxUint64,
	}
	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		var err error
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

	b.Run("vtproto", func(b *testing.B) {
		var err error
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
}

func BenchmarkEncodeInt64(b *testing.B) {
	polyglotData := polyglotBenchmark.I64Data{
		I64: math.MaxInt64,
	}
	polyglotBuf := polyglot.NewBuffer()

	vtData := vtBenchmark.I64Data{
		I64: math.MaxInt64,
	}
	vtBuf := make([]byte, 0, 1024)

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkDecodeInt64(b *testing.B) {
	polyglotData := polyglotBenchmark.I64Data{
		I64: math.MaxInt64,
	}
	polyglotBuf := polyglot.NewBuffer()
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.I64Data{
		I64: math.MaxInt64,
	}
	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		var err error
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

	b.Run("vtproto", func(b *testing.B) {
		var err error
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
}

func BenchmarkEncodeBytes(b *testing.B) {
	randData := make([]byte, 512)
	_, _ = rand.Read(randData)

	polyglotData := polyglotBenchmark.BytesData{
		Bytes: randData,
	}
	polyglotBuf := polyglot.NewBuffer()

	vtData := vtBenchmark.BytesData{
		Bytes: randData,
	}
	vtBuf := make([]byte, 0, 1024)

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		b.SetBytes(512)
		for i := 0; i < b.N; i++ {
			polyglotData.Encode(polyglotBuf)
			polyglotBuf.Reset()
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		b.SetBytes(512)
		for i := 0; i < b.N; i++ {
			_, err = vtData.MarshalToVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
			vtBuf = vtBuf[:0]
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkEncodeBytesParallel(b *testing.B) {
	randData := make([]byte, 512)
	_, _ = rand.Read(randData)

	polyglotData := polyglotBenchmark.BytesData{
		Bytes: randData,
	}

	vtData := vtBenchmark.BytesData{
		Bytes: randData,
	}

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		b.SetBytes(512)
		b.RunParallel(func(pb *testing.PB) {
			polyglotBuf := polyglot.NewBufferSize(1024)
			for pb.Next() {
				polyglotData.Encode(polyglotBuf)
				polyglotBuf.Reset()
			}
			runtime.KeepAlive(polyglotData)
		})
	})

	b.Run("vtproto", func(b *testing.B) {
		b.SetBytes(512)
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

func BenchmarkDecodeBytes(b *testing.B) {
	randData := make([]byte, 512)
	_, _ = rand.Read(randData)

	polyglotData := polyglotBenchmark.BytesData{
		Bytes: randData,
	}
	polyglotBuf := polyglot.NewBuffer()
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.BytesData{
		Bytes: randData,
	}

	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		var err error
		b.SetBytes(512)
		for i := 0; i < b.N; i++ {
			polyglotData.Bytes = nil
			err = polyglotData.Decode(polyglotBytes)
			if err != nil {
				b.Fatal(err)
			}
		}
		runtime.KeepAlive(polyglotData)
	})

	b.Run("vtproto", func(b *testing.B) {
		var err error
		b.SetBytes(512)
		for i := 0; i < b.N; i++ {
			vtData.Reset()
			err = vtData.UnmarshalVT(vtBuf)
			if err != nil {
				b.Fatal(err)
			}
		}
		runtime.KeepAlive(vtData)
	})
}

func BenchmarkDecodeBytesParallel(b *testing.B) {
	randData := make([]byte, 512)
	_, _ = rand.Read(randData)

	polyglotData := polyglotBenchmark.BytesData{
		Bytes: randData,
	}
	polyglotBuf := polyglot.NewBufferSize(1024)
	polyglotData.Encode(polyglotBuf)
	polyglotBytes := polyglotBuf.Bytes()

	vtData := vtBenchmark.BytesData{
		Bytes: randData,
	}

	vtBuf := make([]byte, 1024)
	_, _ = vtData.MarshalToVT(vtBuf)
	vtBuf = vtBuf[:vtData.SizeVT()]

	b.ResetTimer()

	b.Run("polyglot", func(b *testing.B) {
		b.SetBytes(512)
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

	b.Run("vtproto", func(b *testing.B) {
		b.SetBytes(512)
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
