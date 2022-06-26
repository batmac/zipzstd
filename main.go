package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
)

func main() {

	compr := zstd.ZipCompressor(zstd.WithEncoderCRC(true))
	decomp := zstd.ZipDecompressor()

	// Try it out...
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	zw.RegisterCompressor(zstd.ZipMethodWinZip, compr)
	zw.RegisterCompressor(zstd.ZipMethodPKWare, compr)

	// Create 1MB data
	tmp := make([]byte, 1<<20)
	for i := range tmp {
		tmp[i] = byte(i)
	}
	w, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "file1.txt",
		Method: zstd.ZipMethodWinZip,
	})
	if err != nil {
		panic(err)
	}
	w.Write(tmp)

	// Another...
	w, err = zw.CreateHeader(&zip.FileHeader{
		Name:   "file2.txt",
		Method: zstd.ZipMethodPKWare,
	})
	if err != nil {
		panic(err)
	}
	w.Write(tmp)
	zw.Close()

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		panic(err)
	}
	zr.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)
	zr.RegisterDecompressor(zstd.ZipMethodPKWare, decomp)
	for _, file := range zr.File {
		rc, err := file.Open()
		if err != nil {
			panic(err)
		}
		b, err := io.ReadAll(rc)
		if err != nil {
			panic(err)
		}
		if err := rc.Close(); err != nil {
			panic(err)
		}
		if bytes.Equal(b, tmp) {
			fmt.Println(file.Name, "ok")
		} else {
			fmt.Println(file.Name, "mismatch")
		}
	}
}
