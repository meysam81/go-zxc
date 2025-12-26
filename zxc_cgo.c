// Copyright 2025 Meysam Azad
// SPDX-License-Identifier: BSD-3-Clause

// This file includes the ZXC library source files for cgo compilation.
// By placing this in the package directory, cgo will compile it along with
// the Go code, allowing the library to work with 'go get' without a build step.

#include "vendor/zxc/src/lib/zxc_common.c"
#include "vendor/zxc/src/lib/zxc_compress.c"
#include "vendor/zxc/src/lib/zxc_decompress.c"
