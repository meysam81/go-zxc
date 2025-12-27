# go-zxc

[![Go Reference](https://pkg.go.dev/badge/github.com/meysam81/go-zxc.svg)](https://pkg.go.dev/github.com/meysam81/go-zxc)
[![Go Report Card](https://goreportcard.com/badge/github.com/meysam81/go-zxc?style=flat-square)](https://goreportcard.com/report/github.com/meysam81/go-zxc)
[![License](https://img.shields.io/badge/license-Apache--2-blue)](LICENSE)
[![codecov](https://codecov.io/gh/meysam81/go-zxc/graph/badge.svg?token=CJL9KGSWOH)](https://codecov.io/gh/meysam81/go-zxc)
[![GitHub release](https://img.shields.io/github/v/release/meysam81/go-zxc?style=flat-square)](https://github.com/meysam81/go-zxc/releases)
[![GitHub Stars](https://img.shields.io/github/stars/meysam81/go-zxc?style=flat-square&logo=github)](https://github.com/meysam81/go-zxc/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/meysam81/go-zxc?style=flat-square&logo=github)](https://github.com/meysam81/go-zxc/issues)

Go bindings for [ZXC](https://github.com/hellobertrand/zxc), an asymmetric high-performance lossless compression library.

ZXC is designed for "Write Once, Read Many" scenarios where compression speed (build-time) is traded for maximum decompression throughput (run-time). This makes it ideal for content delivery, embedded systems, game assets, firmware, and app bundles.

## Key Features

- **+40% faster decompression** than LZ4 on Apple Silicon
- **+20% faster decompression** than LZ4 on Cloud ARM (Google Axion)
- **Better compression ratios** than LZ4
- **Thread-safe** stateless API suitable for concurrent use
- **Optional checksum** verification for data integrity

## Installation

```bash
go get -u github.com/meysam81/go-zxc
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/meysam81/go-zxc"
)

func main() {
    // Original data
    data := []byte("Hello, ZXC! This is a sample text for compression.")

    // Compress with default options
    compressed, err := zxc.Compress(data, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Original: %d bytes\n", len(data))
    fmt.Printf("Compressed: %d bytes\n", len(compressed))

    // Decompress
    decompressed, err := zxc.Decompress(compressed, len(data), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Decompressed: %s\n", string(decompressed))
}
```

## Compression Levels

ZXC provides four compression levels with different trade-offs:

| Level | Constant        | Description                                          |
| ----- | --------------- | ---------------------------------------------------- |
| 2     | `LevelFast`     | Fastest compression, best for real-time applications |
| 3     | `LevelDefault`  | Recommended: ratio > LZ4, decode speed > LZ4         |
| 4     | `LevelBalanced` | Good ratio and decode speed balance                  |
| 5     | `LevelCompact`  | Highest density, best for storage/firmware/assets    |

```go
// Use compact level for maximum compression
opts := &zxc.Options{
    Level:    zxc.LevelCompact,
    Checksum: true,
}
compressed, err := zxc.Compress(data, opts)
```

## API Reference

### Compression

```go
// Compress data with automatic buffer allocation
compressed, err := zxc.Compress(data, opts)

// Compress into a pre-allocated buffer
n, err := zxc.CompressTo(dst, src, opts)

// Calculate maximum compressed size for pre-allocation
maxSize := zxc.CompressBound(len(data))
```

### Decompression

```go
// Decompress with known original size
decompressed, err := zxc.Decompress(compressed, originalSize, opts)

// Decompress into a pre-allocated buffer
n, err := zxc.DecompressTo(dst, compressed, opts)
```

### Options

```go
type Options struct {
    Level    Level // Compression level (default: LevelDefault)
    Checksum bool  // Enable checksum (default: true)
}
```

### Version

```go
version := zxc.Version() // Returns "0.3.0"
```

## Error Handling

```go
var (
    ErrCompression    // Compression failed
    ErrDecompression  // Decompression failed (corrupted data, checksum mismatch)
    ErrBufferTooSmall // Destination buffer too small
    ErrInvalidLevel   // Invalid compression level
    ErrEmptyInput     // Input data is empty
)
```

## Thread Safety

All functions are thread-safe and can be called concurrently from multiple goroutines. The underlying C library uses a stateless design with caller-allocated buffers.

```go
// Safe for concurrent use
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(data []byte) {
        defer wg.Done()
        compressed, _ := zxc.Compress(data, nil)
        // use compressed...
    }(data)
}
wg.Wait()
```

## Benchmarks

Run benchmarks with:

```bash
go test -bench=. -benchmem
```

Example results on Apple M2:

```shell
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/meysam81/go-zxc
cpu: AMD Ryzen 5 3600 6-Core Processor
```

| Benchmark                           | Iterations | Time per Op   | Throughput   | Allocated  | Allocations |
| ----------------------------------- | ---------- | ------------- | ------------ | ---------- | ----------- |
| BenchmarkStreamCompress-12          | 544        | 2121956 ns/op | 212.07 MB/s  | 400 B/op   | 7 allocs/op |
| BenchmarkStreamDecompress-12        | 524        | 2277227 ns/op | 197.61 MB/s  | 336 B/op   | 6 allocs/op |
| BenchmarkCompress-12                | 138612     | 8934 ns/op    | 5036.80 MB/s | 49152 B/op | 1 allocs/op |
| BenchmarkDecompress-12              | 107698     | 10726 ns/op   | 4195.44 MB/s | 49152 B/op | 1 allocs/op |
| BenchmarkCompressLevels/Fast-12     | 164079     | 8745 ns/op    | 5145.94 MB/s | 49152 B/op | 1 allocs/op |
| BenchmarkCompressLevels/Default-12  | 146947     | 7705 ns/op    | 5840.28 MB/s | 49152 B/op | 1 allocs/op |
| BenchmarkCompressLevels/Balanced-12 | 161617     | 8235 ns/op    | 5464.67 MB/s | 49152 B/op | 1 allocs/op |
| BenchmarkCompressLevels/Compact-12  | 138926     | 8285 ns/op    | 5431.42 MB/s | 49152 B/op | 1 allocs/op |

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.

The underlying ZXC library is Copyright (c) 2025 Bertrand Lebonnois, also under Apache-2.0.
