// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc_test

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/meysam81/go-zxc"
)

func Example() {
	// Original data to compress
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

	// Verify integrity
	if bytes.Equal(data, decompressed) {
		fmt.Println("Data integrity verified!")
	}

	// Output:
	// Original: 50 bytes
	// Compressed: 78 bytes
	// Data integrity verified!
}

func ExampleCompress() {
	data := []byte("Compress this text using ZXC.")

	// Compress with default options (LevelDefault, checksum enabled)
	compressed, err := zxc.Compress(data, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Compressed %d bytes to %d bytes\n", len(data), len(compressed))
}

func ExampleCompress_withOptions() {
	data := []byte("Compress with high compression level for maximum density.")

	// Use compact level for best compression ratio
	opts := &zxc.Options{
		Level:    zxc.LevelCompact,
		Checksum: true,
	}

	compressed, err := zxc.Compress(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Compressed to %d bytes using level %s\n", len(compressed), opts.Level)
}

func ExampleDecompress() {
	original := []byte("This is the original text to be compressed and decompressed.")

	// First compress
	compressed, err := zxc.Compress(original, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Decompress with known original size
	decompressed, err := zxc.Decompress(compressed, len(original), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decompressed: %s\n", string(decompressed))
	// Output:
	// Decompressed: This is the original text to be compressed and decompressed.
}

func ExampleCompressBound() {
	inputSize := 10000

	// Calculate maximum possible compressed size
	maxCompressedSize := zxc.CompressBound(inputSize)

	fmt.Printf("Input size: %d bytes\n", inputSize)
	fmt.Printf("Max compressed size: %d bytes\n", maxCompressedSize)
}

func ExampleCompressTo() {
	data := []byte("Pre-allocate buffer for compression.")

	// Pre-allocate destination buffer
	maxSize := zxc.CompressBound(len(data))
	dst := make([]byte, maxSize)

	// Compress into pre-allocated buffer
	n, err := zxc.CompressTo(dst, data, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Use only the portion that contains compressed data
	compressed := dst[:n]
	fmt.Printf("Wrote %d bytes to pre-allocated buffer\n", len(compressed))
}

func ExampleDecompressTo() {
	original := []byte("Decompress into a pre-allocated buffer.")

	// Compress first
	compressed, err := zxc.Compress(original, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Pre-allocate destination buffer
	dst := make([]byte, len(original))

	// Decompress into buffer
	n, err := zxc.DecompressTo(dst, compressed, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decompressed %d bytes\n", n)
	// Output:
	// Decompressed 39 bytes
}

func ExampleVersion() {
	v := zxc.Version()
	fmt.Printf("ZXC library version: %s\n", v)
}

func ExampleLevel_String() {
	levels := []zxc.Level{
		zxc.LevelFast,
		zxc.LevelDefault,
		zxc.LevelBalanced,
		zxc.LevelCompact,
	}

	for _, level := range levels {
		fmt.Printf("Level %d: %s\n", level, level.String())
	}

	// Output:
	// Level 2: Fast
	// Level 3: Default
	// Level 4: Balanced
	// Level 5: Compact
}

func ExampleStreamCompress() {
	// Create input file with data to compress
	inputFile, err := os.CreateTemp("", "input-*.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(inputFile.Name()) }()
	defer func() { _ = inputFile.Close() }()

	// Write some data
	data := bytes.Repeat([]byte("Example data for streaming compression. "), 100)
	if _, err := inputFile.Write(data); err != nil {
		log.Fatal(err)
	}

	// Seek back to beginning for reading
	if _, err := inputFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Create output file for compressed data
	outputFile, err := os.CreateTemp("", "compressed-*.zxc")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(outputFile.Name()) }()
	defer func() { _ = outputFile.Close() }()

	// Compress using streaming API with multi-threaded pipeline
	opts := &zxc.StreamOptions{
		Level:    zxc.LevelDefault,
		Checksum: true,
		Threads:  0, // auto-detect CPU cores
	}

	compressedBytes, err := zxc.StreamCompress(inputFile, outputFile, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Compressed %d bytes to %d bytes\n", len(data), compressedBytes)

	// Output:
}

func ExampleStreamDecompress() {
	// First, create and compress a file
	inputFile, err := os.CreateTemp("", "input-*.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(inputFile.Name()) }()

	data := bytes.Repeat([]byte("Data for streaming decompression example. "), 100)
	if _, err := inputFile.Write(data); err != nil {
		log.Fatal(err)
	}
	if _, err := inputFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	compressedFile, err := os.CreateTemp("", "compressed-*.zxc")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(compressedFile.Name()) }()

	if _, err := zxc.StreamCompress(inputFile, compressedFile, nil); err != nil {
		log.Fatal(err)
	}
	_ = inputFile.Close()
	_ = compressedFile.Close()

	// Now decompress the file
	compressedFile, err = os.Open(compressedFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = compressedFile.Close() }()

	outputFile, err := os.CreateTemp("", "output-*.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(outputFile.Name()) }()
	defer func() { _ = outputFile.Close() }()

	// Decompress using streaming API
	opts := &zxc.StreamOptions{
		Checksum: true,
		Threads:  4, // use 4 worker threads
	}

	decompressedBytes, err := zxc.StreamDecompress(compressedFile, outputFile, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decompressed %d bytes\n", decompressedBytes)

	// Output:
	// Decompressed 4200 bytes
}
