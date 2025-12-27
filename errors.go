// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: BSD-3-Clause

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
)
