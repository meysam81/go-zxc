# Changelog

## [0.2.2](https://github.com/meysam81/go-zxc/compare/v0.2.1...v0.2.2) (2025-12-28)


### Features

* **CI:** add multiplatform test job ([ae40e75](https://github.com/meysam81/go-zxc/commit/ae40e753ad25e33ffdf77f836d574fedbfd3f67e))


### Bug Fixes

* **CI:** do not fail on tests go versions ([06390d7](https://github.com/meysam81/go-zxc/commit/06390d74920326b3b4039b956aced1859bc211d3))
* **CI:** drop test go v1.20 ([8de8fec](https://github.com/meysam81/go-zxc/commit/8de8fec3c8317ee85a4d4acded8d8fea0167aff1))
* **CI:** quote the versions to avoid interpretation ([c253797](https://github.com/meysam81/go-zxc/commit/c2537972ab241deb4e7d028f30dd109a708d91c2))
* **CI:** remove flto from ldflags ([e766d3d](https://github.com/meysam81/go-zxc/commit/e766d3dada0b5d23cccca0c2887b9129e8a42c86))
* **CI:** remove windows tests ([7e5592c](https://github.com/meysam81/go-zxc/commit/7e5592c160266d03b11fcb885246057dbc86f48c))
* **CI:** specify wildcard for latest go version ([ba72b16](https://github.com/meysam81/go-zxc/commit/ba72b162e1eebb7e35e0adaca9802f870e7658bb))

## [0.2.1](https://github.com/meysam81/go-zxc/compare/v0.2.0...v0.2.1) (2025-12-27)


### Bug Fixes

* **CI:** bring back the root C file ([3538fa2](https://github.com/meysam81/go-zxc/commit/3538fa20fc0ee111be5a8758bea6f7d34079c8d5))
* **CI:** specify ld flags to cgo ([73c6252](https://github.com/meysam81/go-zxc/commit/73c6252a94ae191203da425aea3b488d3b4bcc5c))
* **doc:** add cgo installation note ([89ce5f9](https://github.com/meysam81/go-zxc/commit/89ce5f94a74c9b044c2ae7f6adb69da6738640ed))
* **doc:** update benchmark HW note ([794469e](https://github.com/meysam81/go-zxc/commit/794469e19efeccc1f410b3caae620cb062ed2087))
* include cflags for stream bindings ([cf78618](https://github.com/meysam81/go-zxc/commit/cf7861826e492df3d5e6c6504fc92c4f8b1b799e))
* link to internal vendored C files ([4910a7b](https://github.com/meysam81/go-zxc/commit/4910a7bfc3b34513eb56722db42f9307b0105d41))
* vendor and modify include path ([9fa4f14](https://github.com/meysam81/go-zxc/commit/9fa4f14040673b8495c2f81c0df3027777336163))
* vendor C library for ease of downstream imports ([7b38f5a](https://github.com/meysam81/go-zxc/commit/7b38f5a0dc0d57d9ae119cd833fbc87cb1971d61))

## [0.2.0](https://github.com/meysam81/go-zxc/compare/v0.1.0...v0.2.0) (2025-12-27)


### Features

* add streaming API implementation ([#6](https://github.com/meysam81/go-zxc/issues/6)) ([5c15693](https://github.com/meysam81/go-zxc/commit/5c15693f381daf1152a33cf865302bba68b1830d))

## 0.1.0 (2025-12-27)


### Features

* add codecov badge to readme ([c900c0f](https://github.com/meysam81/go-zxc/commit/c900c0fd63c0ff6d070b30a6cfbd35d53ff9e884))
* add go bindings for zxc compression library ([#2](https://github.com/meysam81/go-zxc/issues/2)) ([6458cc8](https://github.com/meysam81/go-zxc/commit/6458cc81d51b9c2a32d8d76d35786bef6eda134d))
* **CI:** add codecov for coverage report ([caff3a3](https://github.com/meysam81/go-zxc/commit/caff3a341b959201c3265979051e0124a946c710))
* **doc:** add benchmark and lower go version for compatibility ([b85efe3](https://github.com/meysam81/go-zxc/commit/b85efe397e3d260fe32a78bff45389dda172e1c2))
* pin to zxc 0.3.0 ([b9521bd](https://github.com/meysam81/go-zxc/commit/b9521bdbba7ad650af710de07b3ad4d5729ad75d))


### Bug Fixes

* **CI:** recurse submodules ([31ee6a9](https://github.com/meysam81/go-zxc/commit/31ee6a9b4f662fe95408319d5d58764233238208))
* **CI:** remove bun install ([8105094](https://github.com/meysam81/go-zxc/commit/8105094cfc4f65bd938cad7eee54392d28b7de9e))
* use apache 2 license ([5908d97](https://github.com/meysam81/go-zxc/commit/5908d97c40e3546ba31c98cc17d8c60137c33001))
