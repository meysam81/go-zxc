// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc

import (
	"bytes"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStreamCompressDecompressRoundtrip(t *testing.T) {
	testData := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 1000)

	tests := []struct {
		name  string
		level Level
	}{
		{"Fast", LevelFast},
		{"Default", LevelDefault},
		{"Balanced", LevelBalanced},
		{"Compact", LevelCompact},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "input.dat")
			compressedPath := filepath.Join(tmpDir, "compressed.zxc")
			outputPath := filepath.Join(tmpDir, "output.dat")

			// Write test data to input file
			if err := os.WriteFile(inputPath, testData, 0644); err != nil {
				t.Fatalf("Failed to write input file: %v", err)
			}

			// Open input file for compression
			inputFile, err := os.Open(inputPath)
			if err != nil {
				t.Fatalf("Failed to open input file: %v", err)
			}

			// Create compressed file
			compressedFile, err := os.Create(compressedPath)
			if err != nil {
				_ = inputFile.Close()
				t.Fatalf("Failed to create compressed file: %v", err)
			}

			opts := &StreamOptions{
				Level:    tt.level,
				Checksum: true,
			}

			// Compress
			compressedBytes, err := StreamCompress(inputFile, compressedFile, opts)
			if err != nil {
				_ = inputFile.Close()
				_ = compressedFile.Close()
				t.Fatalf("StreamCompress() error: %v", err)
			}

			if compressedBytes <= 0 {
				t.Errorf("StreamCompress() returned %d bytes, expected > 0", compressedBytes)
			}

			// Close files to flush buffers
			_ = inputFile.Close()
			_ = compressedFile.Close()

			// Open compressed file for decompression
			compressedFile, err = os.Open(compressedPath)
			if err != nil {
				t.Fatalf("Failed to open compressed file: %v", err)
			}

			// Create output file
			outputFile, err := os.Create(outputPath)
			if err != nil {
				_ = compressedFile.Close()
				t.Fatalf("Failed to create output file: %v", err)
			}

			// Decompress
			decompressedBytes, err := StreamDecompress(compressedFile, outputFile, opts)
			if err != nil {
				_ = compressedFile.Close()
				_ = outputFile.Close()
				t.Fatalf("StreamDecompress() error: %v", err)
			}

			if decompressedBytes != int64(len(testData)) {
				t.Errorf("StreamDecompress() returned %d bytes, expected %d", decompressedBytes, len(testData))
			}

			// Close and verify
			_ = compressedFile.Close()
			_ = outputFile.Close()

			// Read and verify output
			result, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if !bytes.Equal(testData, result) {
				t.Errorf("Roundtrip failed: got %d bytes, want %d bytes", len(result), len(testData))
			}

			t.Logf("Level %s: %d -> %d bytes (%.1f%%)",
				tt.name, len(testData), compressedBytes,
				float64(compressedBytes)*100/float64(len(testData)))
		})
	}
}

func TestStreamCompressWithThreads(t *testing.T) {
	testData := bytes.Repeat([]byte("Testing multi-threaded compression. "), 5000)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	compressedPath := filepath.Join(tmpDir, "compressed.zxc")
	outputPath := filepath.Join(tmpDir, "output.dat")

	// Write test data
	if err := os.WriteFile(inputPath, testData, 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// Test with different thread counts
	threadCounts := []int{0, 1, 2, 4}

	for _, threads := range threadCounts {
		t.Run("", func(t *testing.T) {
			inputFile, err := os.Open(inputPath)
			if err != nil {
				t.Fatalf("Failed to open input file: %v", err)
			}

			compressedFile, err := os.Create(compressedPath)
			if err != nil {
				_ = inputFile.Close()
				t.Fatalf("Failed to create compressed file: %v", err)
			}

			opts := &StreamOptions{
				Level:    LevelDefault,
				Checksum: true,
				Threads:  threads,
			}

			_, err = StreamCompress(inputFile, compressedFile, opts)
			if err != nil {
				_ = inputFile.Close()
				_ = compressedFile.Close()
				t.Fatalf("StreamCompress() with %d threads error: %v", threads, err)
			}

			_ = inputFile.Close()
			_ = compressedFile.Close()

			// Decompress and verify
			compressedFile, err = os.Open(compressedPath)
			if err != nil {
				t.Fatalf("Failed to open compressed file: %v", err)
			}

			outputFile, err := os.Create(outputPath)
			if err != nil {
				_ = compressedFile.Close()
				t.Fatalf("Failed to create output file: %v", err)
			}

			opts.Threads = threads
			_, err = StreamDecompress(compressedFile, outputFile, opts)
			if err != nil {
				_ = compressedFile.Close()
				_ = outputFile.Close()
				t.Fatalf("StreamDecompress() with %d threads error: %v", threads, err)
			}

			_ = compressedFile.Close()
			_ = outputFile.Close()

			result, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}
			if !bytes.Equal(testData, result) {
				t.Errorf("Roundtrip with %d threads failed", threads)
			}
		})
	}
}

func TestStreamCompressNilFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := os.Create(filepath.Join(tmpDir, "temp.dat"))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = tmpFile.Close() }()

	_, err = StreamCompress(nil, tmpFile, nil)
	if !errors.Is(err, ErrStreamNilFile) {
		t.Errorf("StreamCompress(nil, file) = %v, want ErrStreamNilFile", err)
	}

	_, err = StreamCompress(tmpFile, nil, nil)
	if !errors.Is(err, ErrStreamNilFile) {
		t.Errorf("StreamCompress(file, nil) = %v, want ErrStreamNilFile", err)
	}
}

func TestStreamDecompressNilFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := os.Create(filepath.Join(tmpDir, "temp.dat"))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = tmpFile.Close() }()

	_, err = StreamDecompress(nil, tmpFile, nil)
	if !errors.Is(err, ErrStreamNilFile) {
		t.Errorf("StreamDecompress(nil, file) = %v, want ErrStreamNilFile", err)
	}

	_, err = StreamDecompress(tmpFile, nil, nil)
	if !errors.Is(err, ErrStreamNilFile) {
		t.Errorf("StreamDecompress(file, nil) = %v, want ErrStreamNilFile", err)
	}
}

func TestStreamCompressInvalidLevel(t *testing.T) {
	tmpDir := t.TempDir()

	inputPath := filepath.Join(tmpDir, "input.dat")
	if err := os.WriteFile(inputPath, []byte("test data"), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Failed to open input file: %v", err)
	}
	defer func() { _ = inputFile.Close() }()

	outputFile, err := os.Create(filepath.Join(tmpDir, "output.zxc"))
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer func() { _ = outputFile.Close() }()

	_, err = StreamCompress(inputFile, outputFile, &StreamOptions{Level: Level(100)})
	if !errors.Is(err, ErrInvalidLevel) {
		t.Errorf("StreamCompress with Level(100) = %v, want ErrInvalidLevel", err)
	}
}

func TestStreamCompressWithoutChecksum(t *testing.T) {
	testData := bytes.Repeat([]byte("Testing compression without checksum. "), 500)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	compressedPath := filepath.Join(tmpDir, "compressed.zxc")
	outputPath := filepath.Join(tmpDir, "output.dat")

	if err := os.WriteFile(inputPath, testData, 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Failed to open input file: %v", err)
	}

	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		_ = inputFile.Close()
		t.Fatalf("Failed to create compressed file: %v", err)
	}

	opts := &StreamOptions{
		Level:    LevelDefault,
		Checksum: false,
	}

	_, err = StreamCompress(inputFile, compressedFile, opts)
	if err != nil {
		_ = inputFile.Close()
		_ = compressedFile.Close()
		t.Fatalf("StreamCompress() error: %v", err)
	}

	_ = inputFile.Close()
	_ = compressedFile.Close()

	// Decompress
	compressedFile, err = os.Open(compressedPath)
	if err != nil {
		t.Fatalf("Failed to open compressed file: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		_ = compressedFile.Close()
		t.Fatalf("Failed to create output file: %v", err)
	}

	_, err = StreamDecompress(compressedFile, outputFile, opts)
	if err != nil {
		_ = compressedFile.Close()
		_ = outputFile.Close()
		t.Fatalf("StreamDecompress() error: %v", err)
	}

	_ = compressedFile.Close()
	_ = outputFile.Close()

	result, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if !bytes.Equal(testData, result) {
		t.Error("Roundtrip without checksum failed")
	}
}

func TestStreamCompressNilOptions(t *testing.T) {
	testData := bytes.Repeat([]byte("Testing with nil options. "), 500)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	compressedPath := filepath.Join(tmpDir, "compressed.zxc")
	outputPath := filepath.Join(tmpDir, "output.dat")

	if err := os.WriteFile(inputPath, testData, 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Failed to open input file: %v", err)
	}

	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		_ = inputFile.Close()
		t.Fatalf("Failed to create compressed file: %v", err)
	}

	// nil options should use defaults
	_, err = StreamCompress(inputFile, compressedFile, nil)
	if err != nil {
		_ = inputFile.Close()
		_ = compressedFile.Close()
		t.Fatalf("StreamCompress() with nil opts error: %v", err)
	}

	_ = inputFile.Close()
	_ = compressedFile.Close()

	compressedFile, err = os.Open(compressedPath)
	if err != nil {
		t.Fatalf("Failed to open compressed file: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		_ = compressedFile.Close()
		t.Fatalf("Failed to create output file: %v", err)
	}

	_, err = StreamDecompress(compressedFile, outputFile, nil)
	if err != nil {
		_ = compressedFile.Close()
		_ = outputFile.Close()
		t.Fatalf("StreamDecompress() with nil opts error: %v", err)
	}

	_ = compressedFile.Close()
	_ = outputFile.Close()

	result, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if !bytes.Equal(testData, result) {
		t.Error("Roundtrip with nil options failed")
	}
}

func TestStreamLargeFile(t *testing.T) {
	// Test with 5MB of random data
	size := 5 * 1024 * 1024
	testData := make([]byte, size)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	compressedPath := filepath.Join(tmpDir, "compressed.zxc")
	outputPath := filepath.Join(tmpDir, "output.dat")

	if err := os.WriteFile(inputPath, testData, 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Failed to open input file: %v", err)
	}

	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		_ = inputFile.Close()
		t.Fatalf("Failed to create compressed file: %v", err)
	}

	compressedBytes, err := StreamCompress(inputFile, compressedFile, nil)
	if err != nil {
		_ = inputFile.Close()
		_ = compressedFile.Close()
		t.Fatalf("StreamCompress() error: %v", err)
	}

	t.Logf("Compressed %d bytes to %d bytes (%.1f%%)",
		size, compressedBytes, float64(compressedBytes)*100/float64(size))

	_ = inputFile.Close()
	_ = compressedFile.Close()

	compressedFile, err = os.Open(compressedPath)
	if err != nil {
		t.Fatalf("Failed to open compressed file: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		_ = compressedFile.Close()
		t.Fatalf("Failed to create output file: %v", err)
	}

	decompressedBytes, err := StreamDecompress(compressedFile, outputFile, nil)
	if err != nil {
		_ = compressedFile.Close()
		_ = outputFile.Close()
		t.Fatalf("StreamDecompress() error: %v", err)
	}

	if decompressedBytes != int64(size) {
		t.Errorf("Decompressed %d bytes, expected %d", decompressedBytes, size)
	}

	_ = compressedFile.Close()
	_ = outputFile.Close()

	result, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if !bytes.Equal(testData, result) {
		t.Error("Large file roundtrip failed")
	}
}

func BenchmarkStreamCompress(b *testing.B) {
	data := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 10000)

	tmpDir := b.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	if err := os.WriteFile(inputPath, data, 0644); err != nil {
		b.Fatalf("Failed to write input file: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		inputFile, err := os.Open(inputPath)
		if err != nil {
			b.Fatal(err)
		}
		compressedFile, err := os.Create(filepath.Join(tmpDir, "compressed.zxc"))
		if err != nil {
			_ = inputFile.Close()
			b.Fatal(err)
		}

		_, err = StreamCompress(inputFile, compressedFile, nil)
		if err != nil {
			_ = inputFile.Close()
			_ = compressedFile.Close()
			b.Fatal(err)
		}

		_ = inputFile.Close()
		_ = compressedFile.Close()
	}
}

func BenchmarkStreamDecompress(b *testing.B) {
	data := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 10000)

	tmpDir := b.TempDir()
	inputPath := filepath.Join(tmpDir, "input.dat")
	compressedPath := filepath.Join(tmpDir, "compressed.zxc")
	outputPath := filepath.Join(tmpDir, "output.dat")

	if err := os.WriteFile(inputPath, data, 0644); err != nil {
		b.Fatalf("Failed to write input file: %v", err)
	}

	// Create compressed file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		b.Fatal(err)
	}
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		_ = inputFile.Close()
		b.Fatal(err)
	}
	_, err = StreamCompress(inputFile, compressedFile, nil)
	if err != nil {
		_ = inputFile.Close()
		_ = compressedFile.Close()
		b.Fatal(err)
	}
	_ = inputFile.Close()
	_ = compressedFile.Close()

	b.ResetTimer()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		compressedFile, err := os.Open(compressedPath)
		if err != nil {
			b.Fatal(err)
		}
		outputFile, err := os.Create(outputPath)
		if err != nil {
			_ = compressedFile.Close()
			b.Fatal(err)
		}

		_, err = StreamDecompress(compressedFile, outputFile, nil)
		if err != nil {
			_ = compressedFile.Close()
			_ = outputFile.Close()
			b.Fatal(err)
		}

		_ = compressedFile.Close()
		_ = outputFile.Close()
	}
}
