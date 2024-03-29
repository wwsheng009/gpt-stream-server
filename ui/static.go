package ui

import (
	"net/http"
	"os"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: root, Fallback: "index.html"}
	return &binaryFileSystem{
		fs,
	}
}
