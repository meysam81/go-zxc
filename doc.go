// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

// Package zxc provides Go bindings for the ZXC asymmetric high-performance
// lossless compression library.
//
// ZXC is designed for "Write Once, Read Many" scenarios where compression
// speed (build-time) is traded for maximum decompression throughput (run-time).
// This makes it ideal for content delivery, embedded systems, game assets,
// firmware, and app bundles.
//
// # Key Features
//
//   - Outperforms LZ4 decompression by +40% on Apple Silicon and +20% on Cloud ARM
//   - Better compression ratios than LZ4
//   - Thread-safe, stateless API suitable for concurrent use
//   - Optional checksum verification for data integrity
//   - Streaming API for large files with multi-threaded processing
//
// # Basic Usage
//
// Compress data using [Compress] and decompress using [Decompress]:
//
//	compressed, err := zxc.Compress(data, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	decompressed, err := zxc.Decompress(compressed, len(data), nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Compression Levels
//
// ZXC provides four compression levels with different trade-offs:
//
//   - [LevelFast]: Fastest compression, best for real-time applications
//   - [LevelDefault]: Recommended for most use cases (ratio > LZ4, decode speed > LZ4)
//   - [LevelBalanced]: Good ratio and decode speed balance
//   - [LevelCompact]: Highest density, best for storage/firmware/assets
//
// # Options
//
// Use [Options] to customize compression behavior:
//
//	opts := &zxc.Options{
//	    Level:    zxc.LevelCompact,
//	    Checksum: true,
//	}
//	compressed, err := zxc.Compress(data, opts)
//
// # Buffer Management
//
// For applications that need to control memory allocation, use [CompressBound]
// to calculate the maximum compressed size and [CompressTo]/[DecompressTo] to
// write to pre-allocated buffers:
//
//	maxSize := zxc.CompressBound(len(data))
//	dst := make([]byte, maxSize)
//
//	n, err := zxc.CompressTo(dst, data, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	compressed := dst[:n]
//
// # Streaming API
//
// For large files that may not fit in memory, use the streaming API with
// [StreamCompress] and [StreamDecompress]:
//
//	inputFile, _ := os.Open("large-file.dat")
//	outputFile, _ := os.Create("large-file.dat.zxc")
//
//	opts := &zxc.StreamOptions{
//	    Level:    zxc.LevelDefault,
//	    Checksum: true,
//	    Threads:  0,
//	}
//
//	compressedBytes, err := zxc.StreamCompress(inputFile, outputFile, opts)
//
// The streaming API uses a multi-threaded pipeline with asynchronous I/O to
// maximize throughput. Set Threads to 0 to auto-detect the number of CPU cores.
//
// # Thread Safety
//
// All functions in this package are thread-safe and can be called concurrently
// from multiple goroutines. The underlying C library uses a stateless design
// with caller-allocated buffers.
//
// # Error Handling
//
// Functions return [ErrCompression] or [ErrDecompression] when the operation
// fails. [ErrBufferTooSmall] is returned when the destination buffer cannot
// hold the result. Streaming functions return [ErrStreamCompression],
// [ErrStreamDecompression], [ErrStreamNilFile], or [ErrStreamOpen] for
// streaming-specific errors.
package zxc
