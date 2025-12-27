// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: BSD-3-Clause

package zxc

// Level represents a ZXC compression level.
// Higher levels produce smaller output at the cost of slower compression.
// Decompression speed is generally consistent across all levels.
type Level int

const (
	// LevelFast provides the fastest compression speed.
	// Best suited for real-time applications where compression time is critical.
	// Compression level: 2
	LevelFast Level = 2

	// LevelDefault is the recommended compression level for most use cases.
	// Provides better compression ratio than LZ4 with faster decompression speed.
	// Compression level: 3
	LevelDefault Level = 3

	// LevelBalanced offers a good trade-off between compression ratio and speed.
	// Suitable for applications that need both reasonable compression time and good ratios.
	// Compression level: 4
	LevelBalanced Level = 4

	// LevelCompact provides the highest compression density.
	// Best for storage-constrained environments like firmware, embedded systems, and assets.
	// Compression level: 5
	LevelCompact Level = 5
)

// String returns a human-readable name for the compression level.
func (l Level) String() string {
	switch l {
	case LevelFast:
		return "Fast"
	case LevelDefault:
		return "Default"
	case LevelBalanced:
		return "Balanced"
	case LevelCompact:
		return "Compact"
	default:
		return "Unknown"
	}
}

// Valid reports whether the level is a valid ZXC compression level.
func (l Level) Valid() bool {
	return l >= LevelFast && l <= LevelCompact
}
