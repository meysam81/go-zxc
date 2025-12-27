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
//
// # Basic Usage
//
// Compress data using [Compress] and decompress using [Decompress]:
//
//	// Compress with default level
//	compressed, err := zxc.Compress(data, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Decompress
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
//	// Pre-allocate destination buffer
//	maxSize := zxc.CompressBound(len(data))
//	dst := make([]byte, maxSize)
//
//	// Compress into the buffer
//	n, err := zxc.CompressTo(dst, data, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	compressed := dst[:n]
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
// hold the result.
package zxc
