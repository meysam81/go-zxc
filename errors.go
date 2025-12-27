// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc

import "errors"

var (
	// ErrCompression is returned when compression fails.
	// This can occur due to invalid input data or internal compression errors.
	ErrCompression = errors.New("zxc: compression failed")

	// ErrDecompression is returned when decompression fails.
	// This can occur due to corrupted data, invalid headers, or checksum mismatches.
	ErrDecompression = errors.New("zxc: decompression failed")

	// ErrBufferTooSmall is returned when the destination buffer is too small
	// to hold the compressed or decompressed data.
	ErrBufferTooSmall = errors.New("zxc: destination buffer too small")

	// ErrInvalidLevel is returned when an invalid compression level is specified.
	ErrInvalidLevel = errors.New("zxc: invalid compression level")

	// ErrEmptyInput is returned when the input data is empty.
	ErrEmptyInput = errors.New("zxc: input data is empty")

	// ErrStreamCompression is returned when streaming compression fails.
	// This can occur due to I/O errors, memory allocation failures, or
	// thread synchronization issues.
	ErrStreamCompression = errors.New("zxc: stream compression failed")

	// ErrStreamDecompression is returned when streaming decompression fails.
	// This can occur due to corrupted data, I/O errors, or checksum mismatches.
	ErrStreamDecompression = errors.New("zxc: stream decompression failed")

	// ErrStreamNilFile is returned when a nil file is passed to streaming functions.
	ErrStreamNilFile = errors.New("zxc: nil file provided")

	// ErrStreamOpen is returned when opening a file stream fails.
	ErrStreamOpen = errors.New("zxc: failed to open file stream")
)
