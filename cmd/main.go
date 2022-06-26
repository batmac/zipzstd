package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"archive/zip"

	"github.com/klauspost/compress/zstd"
	flag "github.com/spf13/pflag"
)

var (
	argDecompress  = flag.BoolP("decompress", "x", false, "decompress files")
	argArchiveFile = flag.StringP("file", "f", "-", "file")
)

func main() {
	flag.Parse()

	if *argDecompress {
		log.Println("decompress")
	}

	archiveFile := *argArchiveFile
	files := flag.Args()[1:]

	if len(files) == 0 {
		panic("no files")
	}

	compr := zstd.ZipCompressor(zstd.WithEncoderCRC(true))
	decomp := zstd.ZipDecompressor()

	// Try it out...
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	zw.RegisterCompressor(zstd.ZipMethodWinZip, compr)
	zw.RegisterCompressor(zstd.ZipMethodPKWare, compr)

	for _, path := range files {
		// Create 1MB data
		tmp := make([]byte, 1<<20)
		for i := range tmp {
			tmp[i] = byte(i)
		}
		w, err := zw.CreateHeader(&zip.FileHeader{
			Name:   path,
			Method: zstd.ZipMethodWinZip,
		})
		if err != nil {
			panic(err)
		}
		w.Write(tmp)
	}

	zw.Close()

	f, err := os.Create(archiveFile)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(f, &buf)
	if err != nil {
		panic(err)
	}
	log.Printf("wrote %v bytes.", n)
	if f.Close(); err != nil {
		panic(err)
	}

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
