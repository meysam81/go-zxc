// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc

/*
#cgo CFLAGS: -I${SRCDIR}/internal/czxc -O3 -flto -fomit-frame-pointer -fstrict-aliasing -ffast-math -pthread -DNDEBUG
#cgo linux CFLAGS: -D_GNU_SOURCE
#cgo linux LDFLAGS: -flto -pthread -lm
#cgo darwin LDFLAGS: -lm

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "zxc.h"

// Helper to create FILE* from file descriptor with dup to avoid double-close issues
static FILE* zxc_fdopen_read(int fd) {
    int duped = dup(fd);
    if (duped < 0) return NULL;
    FILE *f = fdopen(duped, "rb");
    if (f == NULL) {
        close(duped);
        return NULL;
    }
    return f;
}

static FILE* zxc_fdopen_write(int fd) {
    int duped = dup(fd);
    if (duped < 0) return NULL;
    FILE *f = fdopen(duped, "wb");
    if (f == NULL) {
        close(duped);
        return NULL;
    }
    return f;
}
*/
import "C"

import (
	"os"
	"runtime"
)

// StreamOptions configures streaming compression/decompression behavior.
type StreamOptions struct {
	// Level specifies the compression level.
	// If zero, LevelDefault is used.
	// Only used for compression; ignored during decompression.
	Level Level

	// Checksum enables checksum calculation during compression
	// and verification during decompression.
	Checksum bool

	// Threads specifies the number of worker threads to spawn.
	// If zero, auto-detects the number of CPU cores.
	Threads int
}

// streamDefaults returns a copy of the options with default values applied.
func (o *StreamOptions) streamDefaults() StreamOptions {
	if o == nil {
		return StreamOptions{
			Level:    LevelDefault,
			Checksum: true,
			Threads:  0, // auto-detect
		}
	}
	opts := *o
	if opts.Level == 0 {
		opts.Level = LevelDefault
	}
	if opts.Threads < 0 {
		opts.Threads = 0
	}
	return opts
}

// StreamCompress compresses data from the input file to the output file using
// a multi-threaded streaming pipeline.
//
// The streaming API is designed for large files that may not fit in memory.
// It uses an asynchronous producer-consumer architecture with a ring buffer
// to maximize throughput by parallelizing I/O and CPU-bound compression.
//
// If opts is nil, default options are used (LevelDefault, checksum enabled,
// auto-detected thread count).
//
// Returns the total number of compressed bytes written, or an error if
// compression fails.
//
// Note: Both input and output files must be open and ready for reading/writing.
// The caller is responsible for closing the files after this function returns.
func StreamCompress(input, output *os.File, opts *StreamOptions) (int64, error) {
	if input == nil || output == nil {
		return 0, ErrStreamNilFile
	}

	o := opts.streamDefaults()
	if !o.Level.Valid() {
		return 0, ErrInvalidLevel
	}

	checksumFlag := 0
	if o.Checksum {
		checksumFlag = 1
	}

	// Get file descriptors
	inFd := C.int(input.Fd())
	outFd := C.int(output.Fd())

	// Create FILE* from file descriptors
	fIn := C.zxc_fdopen_read(inFd)
	if fIn == nil {
		return 0, ErrStreamOpen
	}
	defer C.fclose(fIn)

	fOut := C.zxc_fdopen_write(outFd)
	if fOut == nil {
		return 0, ErrStreamOpen
	}
	defer C.fclose(fOut)

	// Pin to OS thread while doing the streaming operation
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	result := C.zxc_stream_compress(
		fIn,
		fOut,
		C.int(o.Threads),
		C.int(o.Level),
		C.int(checksumFlag),
	)

	if result < 0 {
		return 0, ErrStreamCompression
	}

	return int64(result), nil
}

// StreamDecompress decompresses data from the input file to the output file
// using a multi-threaded streaming pipeline.
//
// The streaming API is designed for large files that may not fit in memory.
// It uses an asynchronous producer-consumer architecture with a ring buffer
// to maximize throughput by parallelizing I/O and CPU-bound decompression.
//
// If opts is nil, default options are used (checksum verification enabled,
// auto-detected thread count).
//
// Returns the total number of decompressed bytes written, or an error if
// decompression fails.
//
// Note: Both input and output files must be open and ready for reading/writing.
// The caller is responsible for closing the files after this function returns.
func StreamDecompress(input, output *os.File, opts *StreamOptions) (int64, error) {
	if input == nil || output == nil {
		return 0, ErrStreamNilFile
	}

	o := opts.streamDefaults()

	checksumFlag := 0
	if o.Checksum {
		checksumFlag = 1
	}

	// Get file descriptors
	inFd := C.int(input.Fd())
	outFd := C.int(output.Fd())

	// Create FILE* from file descriptors
	fIn := C.zxc_fdopen_read(inFd)
	if fIn == nil {
		return 0, ErrStreamOpen
	}
	defer C.fclose(fIn)

	fOut := C.zxc_fdopen_write(outFd)
	if fOut == nil {
		return 0, ErrStreamOpen
	}
	defer C.fclose(fOut)

	// Pin to OS thread while doing the streaming operation
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	result := C.zxc_stream_decompress(
		fIn,
		fOut,
		C.int(o.Threads),
		C.int(checksumFlag),
	)

	if result < 0 {
		return 0, ErrStreamDecompression
	}

	return int64(result), nil
}
