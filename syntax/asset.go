package syntax

// @TODO: DEPRECIDADO

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

// AssetFileInfo is a fs.FileInfo
type AssetFileInfo struct {
	name string
	time time.Time
	size int64
}

func (i AssetFileInfo) Name() string       { return i.name }
func (i AssetFileInfo) Size() int64        { return i.size }
func (i AssetFileInfo) ModTime() time.Time { return i.time }
func (i AssetFileInfo) Mode() os.FileMode  { return 0444 } // Read for all
func (i AssetFileInfo) IsDir() bool        { return false }
func (i AssetFileInfo) Sys() interface{}   { return nil }

// AssetFile is a http.File
type AssetFile struct {
	*bytes.Reader
	info AssetFileInfo
}

func (f *AssetFile) Stat() (fs.FileInfo, error)               { return f.info, nil }
func (f *AssetFile) Readdir(count int) ([]os.FileInfo, error) { return nil, nil }
func (f *AssetFile) Close() error                             { return nil }

func NewHttpFile(name string, modification time.Time, data []byte) http.File {
	mf := &AssetFile{
		Reader: bytes.NewReader(data),
		info: AssetFileInfo{
			name: name,
			time: modification,
			size: int64(len(data)),
		},
	}

	var f http.File = mf
	return f
}

// SingleFileFileSystem é um http.FileSystem que sempre entrega o mesmo arquivo no método Open
type SingleFileFileSystem struct {
	file http.File
}

func (f SingleFileFileSystem) Open(string) (http.File, error) { return f.file, nil }

func isAsset(name string) bool {
	return strings.HasSuffix(name, ".js") || strings.HasSuffix(name, ".css")
}
