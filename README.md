# About

[![Build status](https://github.com/rgl/go-diskfs-create-fat32-img/workflows/Build/badge.svg)](https://github.com/rgl/go-diskfs-create-fat32-img/actions?query=workflow%3ABuild)

This builds a FAT32 disk image using [go-diskfs](https://github.com/diskfs/go-diskfs).

# Usage (Native)

Build and execute:

```bash
go build
dd if=/dev/urandom of=Setup.bin bs=123456 count=1 && rm -f Setup.bin.img
./go-diskfs-create-fat32-img
```

# Usage (WebAssembly/wasmtime)

**NB** See https://github.com/diskfs/go-diskfs/issues/207.

Build and execute with [wasmtime](https://github.com/bytecodealliance/wasmtime):

```bash
GOOS=wasip1 GOARCH=wasm go build -o go-diskfs-create-fat32-img.wasm
dd if=/dev/urandom of=Setup.bin bs=123456 count=1 && rm -f Setup.bin.img
wasmtime run --dir .::/ go-diskfs-create-fat32-img.wasm
```

# Usage (WebAssembly/wasmer)

**NB** See https://github.com/diskfs/go-diskfs/issues/207.

Build and execute with [wasmer](https://github.com/wasmerio/wasmer):

```bash
GOOS=wasip1 GOARCH=wasm go build -o go-diskfs-create-fat32-img.wasm
dd if=/dev/urandom of=Setup.bin bs=123456 count=1 && rm -f Setup.bin.img
wasmer run --mapdir .::/ go-diskfs-create-fat32-img.wasm
```

**NB** This currently fails with the error (see https://github.com/wasmerio/wasmer/issues/4384):

```
2024/01/03 07:19:31 open Setup.bin: Capabilities insufficient
panic: open Setup.bin: Capabilities insufficient

goroutine 1 [running]:
log.Panic({0x1461f00, 0x1, 0x1})
        /opt/go/src/log/log.go:432 +0x5
main.main()
        /home/vagrant/Projects/go-diskfs-create-fat32-img/main.go:17 +0x6
```
