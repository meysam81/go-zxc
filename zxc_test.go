// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: Apache-2.0

package zxc

import (
	"bytes"
	"crypto/rand"
	"errors"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	if v == "" {
		t.Error("Version() returned empty string")
	}
	// Check format is like "X.Y.Z"
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		t.Errorf("Version() = %q, expected format X.Y.Z", v)
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelFast, "Fast"},
		{LevelDefault, "Default"},
		{LevelBalanced, "Balanced"},
		{LevelCompact, "Compact"},
		{Level(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.level.String()
		if got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestLevelValid(t *testing.T) {
	tests := []struct {
		level Level
		valid bool
	}{
		{LevelFast, true},
		{LevelDefault, true},
		{LevelBalanced, true},
		{LevelCompact, true},
		{Level(0), false},
		{Level(1), false},
		{Level(6), false},
		{Level(100), false},
	}

	for _, tt := range tests {
		got := tt.level.Valid()
		if got != tt.valid {
			t.Errorf("Level(%d).Valid() = %v, want %v", tt.level, got, tt.valid)
		}
	}
}

func TestCompressBound(t *testing.T) {
	tests := []int{0, 1, 100, 1000, 10000, 100000}

	for _, size := range tests {
		bound := CompressBound(size)
		if bound < size {
			t.Errorf("CompressBound(%d) = %d, expected >= %d", size, bound, size)
		}
	}
}

func TestCompressDecompressRoundtrip(t *testing.T) {
	testData := []byte("Hello, ZXC! This is a test message for compression and decompression. " +
		"The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs.")

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
			opts := &Options{
				Level:    tt.level,
				Checksum: true,
			}

			// Compress
			compressed, err := Compress(testData, opts)
			if err != nil {
				t.Fatalf("Compress() error: %v", err)
			}

			if len(compressed) == 0 {
				t.Fatal("Compress() returned empty data")
			}

			// Decompress
			decompressed, err := Decompress(compressed, len(testData), opts)
			if err != nil {
				t.Fatalf("Decompress() error: %v", err)
			}

			// Verify
			if !bytes.Equal(testData, decompressed) {
				t.Errorf("Roundtrip failed: got %d bytes, want %d bytes", len(decompressed), len(testData))
			}
		})
	}
}

func TestCompressDecompressWithoutChecksum(t *testing.T) {
	testData := []byte("Testing compression without checksum verification.")

	opts := &Options{
		Level:    LevelDefault,
		Checksum: false,
	}

	compressed, err := Compress(testData, opts)
	if err != nil {
		t.Fatalf("Compress() error: %v", err)
	}

	decompressed, err := Decompress(compressed, len(testData), opts)
	if err != nil {
		t.Fatalf("Decompress() error: %v", err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("Roundtrip without checksum failed")
	}
}

func TestCompressNilOptions(t *testing.T) {
	testData := []byte("Testing with nil options.")

	compressed, err := Compress(testData, nil)
	if err != nil {
		t.Fatalf("Compress() with nil opts error: %v", err)
	}

	decompressed, err := Decompress(compressed, len(testData), nil)
	if err != nil {
		t.Fatalf("Decompress() with nil opts error: %v", err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("Roundtrip with nil options failed")
	}
}

func TestCompressEmptyInput(t *testing.T) {
	_, err := Compress([]byte{}, nil)
	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Compress(empty) = %v, want ErrEmptyInput", err)
	}

	_, err = Compress(nil, nil)
	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Compress(nil) = %v, want ErrEmptyInput", err)
	}
}

func TestDecompressEmptyInput(t *testing.T) {
	_, err := Decompress([]byte{}, 100, nil)
	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Decompress(empty) = %v, want ErrEmptyInput", err)
	}
}

func TestCompressInvalidLevel(t *testing.T) {
	testData := []byte("Test data")

	// Level(0) is treated as "use default", not an error
	_, err := Compress(testData, &Options{Level: Level(0)})
	if err != nil {
		t.Errorf("Compress with Level(0) should use default, got error: %v", err)
	}

	// Level(1) is below the valid range (2-5)
	_, err = Compress(testData, &Options{Level: Level(1)})
	if !errors.Is(err, ErrInvalidLevel) {
		t.Errorf("Compress with Level(1) = %v, want ErrInvalidLevel", err)
	}

	// Level(100) is above the valid range
	_, err = Compress(testData, &Options{Level: Level(100)})
	if !errors.Is(err, ErrInvalidLevel) {
		t.Errorf("Compress with Level(100) = %v, want ErrInvalidLevel", err)
	}
}

func TestCompressTo(t *testing.T) {
	testData := []byte("Test data for CompressTo function.")

	maxSize := CompressBound(len(testData))
	dst := make([]byte, maxSize)

	n, err := CompressTo(dst, testData, nil)
	if err != nil {
		t.Fatalf("CompressTo() error: %v", err)
	}

	if n <= 0 || n > maxSize {
		t.Errorf("CompressTo() returned n=%d, expected 0 < n <= %d", n, maxSize)
	}

	// Verify we can decompress
	decompressed, err := Decompress(dst[:n], len(testData), nil)
	if err != nil {
		t.Fatalf("Decompress() error: %v", err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("CompressTo roundtrip failed")
	}
}

func TestDecompressTo(t *testing.T) {
	testData := []byte("Test data for DecompressTo function.")

	compressed, err := Compress(testData, nil)
	if err != nil {
		t.Fatalf("Compress() error: %v", err)
	}

	dst := make([]byte, len(testData)+100) // Extra space

	n, err := DecompressTo(dst, compressed, nil)
	if err != nil {
		t.Fatalf("DecompressTo() error: %v", err)
	}

	if n != len(testData) {
		t.Errorf("DecompressTo() returned n=%d, expected %d", n, len(testData))
	}

	if !bytes.Equal(testData, dst[:n]) {
		t.Error("DecompressTo roundtrip failed")
	}
}

func TestLargeData(t *testing.T) {
	// Test with 1MB of random data
	size := 1024 * 1024
	testData := make([]byte, size)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	compressed, err := Compress(testData, nil)
	if err != nil {
		t.Fatalf("Compress() error: %v", err)
	}

	t.Logf("Compressed %d bytes to %d bytes (%.1f%%)",
		len(testData), len(compressed), float64(len(compressed))*100/float64(len(testData)))

	decompressed, err := Decompress(compressed, size, nil)
	if err != nil {
		t.Fatalf("Decompress() error: %v", err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("Large data roundtrip failed")
	}
}

func TestRepetitiveData(t *testing.T) {
	// Repetitive data should compress very well
	testData := bytes.Repeat([]byte("ABCDEFGHIJ"), 10000)

	compressed, err := Compress(testData, nil)
	if err != nil {
		t.Fatalf("Compress() error: %v", err)
	}

	ratio := float64(len(compressed)) * 100 / float64(len(testData))
	t.Logf("Compressed %d bytes to %d bytes (%.1f%%)", len(testData), len(compressed), ratio)

	// Repetitive data should compress to less than 50% of original
	if ratio > 50 {
		t.Errorf("Compression ratio %.1f%% is worse than expected for repetitive data", ratio)
	}

	decompressed, err := Decompress(compressed, len(testData), nil)
	if err != nil {
		t.Fatalf("Decompress() error: %v", err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("Repetitive data roundtrip failed")
	}
}

func BenchmarkCompress(b *testing.B) {
	data := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 1000)

	b.ResetTimer()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		_, err := Compress(data, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecompress(b *testing.B) {
	data := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 1000)
	compressed, err := Compress(data, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		_, err := Decompress(compressed, len(data), nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompressLevels(b *testing.B) {
	data := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 1000)

	levels := []struct {
		name  string
		level Level
	}{
		{"Fast", LevelFast},
		{"Default", LevelDefault},
		{"Balanced", LevelBalanced},
		{"Compact", LevelCompact},
	}

	for _, l := range levels {
		b.Run(l.name, func(b *testing.B) {
			opts := &Options{Level: l.level}
			b.SetBytes(int64(len(data)))

			for i := 0; i < b.N; i++ {
				_, err := Compress(data, opts)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
