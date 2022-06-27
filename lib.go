package zipstd

import (
	"archive/zip"

	"github.com/klauspost/compress/zstd"
)

const (
	// re-export, prefer the first one
	ZipMethodWinZip = zstd.ZipMethodWinZip
	ZipMethodPKWare = zstd.ZipMethodPKWare
)

var (
	Compr  = zstd.ZipCompressor(zstd.WithEncoderCRC(false))
	Decomp = zstd.ZipDecompressor()
)

type checkCompressor interface {
	RegisterCompressor(method uint16, comp zip.Compressor)
}
type checkDecompressor interface {
	RegisterDecompressor(method uint16, comp zip.Decompressor)
}

func RegisterDecompressor(z checkDecompressor) {
	z.RegisterDecompressor(zstd.ZipMethodWinZip, Decomp)
	z.RegisterDecompressor(zstd.ZipMethodPKWare, Decomp)
}
func RegisterCompressor(z checkCompressor) {
	z.RegisterCompressor(zstd.ZipMethodWinZip, Compr)
	z.RegisterCompressor(zstd.ZipMethodPKWare, Compr)
}

func Register(z any) {
	if zp, ok := z.(checkCompressor); ok {
		RegisterCompressor(zp)
	}
	if zp, ok := z.(checkDecompressor); ok {
		RegisterDecompressor(zp)
	}
}
