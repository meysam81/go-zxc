// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc

/*
#cgo CFLAGS: -I${SRCDIR}/internal/czxc -O3 -DNDEBUG
#cgo LDFLAGS: -lm -lpthread

#include "zxc.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Version constants for the ZXC library.
const (
	VersionMajor = int(C.ZXC_VERSION_MAJOR)
	VersionMinor = int(C.ZXC_VERSION_MINOR)
	VersionPatch = int(C.ZXC_VERSION_PATCH)
)

// Version returns the ZXC library version as a string in the format "major.minor.patch".
func Version() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}

// Options configures compression behavior.
type Options struct {
	// Level specifies the compression level.
	// If zero, LevelDefault is used.
	Level Level

	// Checksum enables checksum calculation during compression
	// and verification during decompression.
	// When enabled, a checksum is stored with the compressed data
	// and verified upon decompression.
	Checksum bool
}

// defaults returns a copy of the options with default values applied.
func (o *Options) defaults() Options {
	if o == nil {
		return Options{
			Level:    LevelDefault,
			Checksum: true,
		}
	}
	opts := *o
	if opts.Level == 0 {
		opts.Level = LevelDefault
	}
	return opts
}

// CompressBound returns the maximum compressed size for the given input size.
// This value can be used to pre-allocate a buffer that is guaranteed to be
// large enough to hold the compressed output.
func CompressBound(srcSize int) int {
	return int(C.zxc_compress_bound(C.size_t(srcSize)))
}

// Compress compresses the input data and returns the compressed bytes.
// It allocates a new buffer for the output.
//
// If opts is nil, default options are used (LevelDefault with checksum enabled).
//
// Returns an error if compression fails or if the input is empty.
func Compress(src []byte, opts *Options) ([]byte, error) {
	if len(src) == 0 {
		return nil, ErrEmptyInput
	}

	o := opts.defaults()
	if !o.Level.Valid() {
		return nil, ErrInvalidLevel
	}

	// Allocate destination buffer with maximum possible size
	maxSize := CompressBound(len(src))
	dst := make([]byte, maxSize)

	n, err := CompressTo(dst, src, opts)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

// CompressTo compresses src into dst and returns the number of bytes written.
// The destination buffer must be large enough to hold the compressed data.
// Use CompressBound to determine the required size.
//
// If opts is nil, default options are used (LevelDefault with checksum enabled).
//
// Returns ErrBufferTooSmall if dst is too small, or ErrCompression if compression fails.
func CompressTo(dst, src []byte, opts *Options) (int, error) {
	if len(src) == 0 {
		return 0, ErrEmptyInput
	}

	o := opts.defaults()
	if !o.Level.Valid() {
		return 0, ErrInvalidLevel
	}

	checksumFlag := 0
	if o.Checksum {
		checksumFlag = 1
	}

	var srcPtr, dstPtr unsafe.Pointer
	if len(src) > 0 {
		srcPtr = unsafe.Pointer(&src[0])
	}
	if len(dst) > 0 {
		dstPtr = unsafe.Pointer(&dst[0])
	}

	n := C.zxc_compress(
		srcPtr,
		C.size_t(len(src)),
		dstPtr,
		C.size_t(len(dst)),
		C.int(o.Level),
		C.int(checksumFlag),
	)

	if n == 0 {
		// Check if buffer was too small
		if len(dst) < CompressBound(len(src)) {
			return 0, ErrBufferTooSmall
		}
		return 0, ErrCompression
	}

	return int(n), nil
}

// Decompress decompresses the input data and returns the decompressed bytes.
// The expectedSize parameter specifies the expected size of the decompressed data.
// If expectedSize is 0 or negative, a reasonable default is used (10x the compressed size).
//
// If opts is nil, default options are used (checksum verification enabled).
//
// Returns an error if decompression fails, data is corrupted, or checksum validation fails.
func Decompress(src []byte, expectedSize int, opts *Options) ([]byte, error) {
	if len(src) == 0 {
		return nil, ErrEmptyInput
	}

	// If no expected size provided, estimate based on compression ratio
	// ZXC typically achieves 40-60% compression, so 10x is a safe upper bound
	if expectedSize <= 0 {
		expectedSize = len(src) * 10
	}

	dst := make([]byte, expectedSize)

	n, err := DecompressTo(dst, src, opts)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

// DecompressTo decompresses src into dst and returns the number of bytes written.
// The destination buffer must be large enough to hold the decompressed data.
//
// If opts is nil, default options are used (checksum verification enabled).
//
// Returns ErrBufferTooSmall if dst is too small, or ErrDecompression if decompression fails.
func DecompressTo(dst, src []byte, opts *Options) (int, error) {
	if len(src) == 0 {
		return 0, ErrEmptyInput
	}

	o := opts.defaults()

	checksumFlag := 0
	if o.Checksum {
		checksumFlag = 1
	}

	var srcPtr, dstPtr unsafe.Pointer
	if len(src) > 0 {
		srcPtr = unsafe.Pointer(&src[0])
	}
	if len(dst) > 0 {
		dstPtr = unsafe.Pointer(&dst[0])
	}

	n := C.zxc_decompress(
		srcPtr,
		C.size_t(len(src)),
		dstPtr,
		C.size_t(len(dst)),
		C.int(checksumFlag),
	)

	if n == 0 {
		return 0, ErrDecompression
	}

	return int(n), nil
}
