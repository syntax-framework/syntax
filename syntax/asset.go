package syntax

import (
	"bytes"
	"github.com/syntax-framework/shtml/cmn"
	"github.com/syntax-framework/shtml/sht"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

// Bundler responsible for grouping the assets used by the pages and, from that, defining the most optimized way to
// group these resources in order to maximize performance.
type Bundler struct {
	dirty             bool                    // Indica que houve mudança na
	assetByName       map[string]*cmn.Asset   // facilita busca
	assetByPage       map[string][]*cmn.Asset // Lista original de assets por página
	bundleByPageBuild map[string][]*cmn.Asset // Lista processada de assets por página
}

// SetPageAssets define os assets que podem ser consumidos por uma página
func (b *Bundler) SetPageAssets(page string, assets []*cmn.Asset) {
	if b.assetByPage == nil {
		b.assetByPage = map[string][]*cmn.Asset{}
	}
	if len(b.assetByPage[page]) != len(assets) {
		b.dirty = true
		b.assetByPage[page] = assets
	} else {
		// check if assets has changed
		tmp := map[*cmn.Asset]bool{}
		for _, asset := range b.assetByPage[page] {
			tmp[asset] = true
		}
		for _, asset := range assets {
			if tmp[asset] != true {
				b.dirty = true
				break
			}
		}

		if b.dirty {
			b.assetByPage[page] = assets
		}
	}
}

// GetAssets returns all assets that should be displayed on a page
func (b *Bundler) GetAssets(page string) []*cmn.Asset {
	if b.dirty {
		b.build()
		b.dirty = false
	}
	var out []*cmn.Asset
	if bundle, exists := b.bundleByPageBuild[page]; exists {
		for _, asset := range bundle {
			out = append(out, asset)
		}
	}
	return out
}

// GetScripts returns all assets that should be displayed on a page
func (b *Bundler) GetScripts(page string) string {
	buf := &bytes.Buffer{}
	for _, asset := range b.GetAssets(page) {
		// @TODO: meta data
		// defer="defer"
		// crossorigin="anonymous"
		// integrity="sha512-y9xS8icoY1YJrM9plRcX5Ko3dxi37Poz/u9gxSWevG/YT1NTcKKzgSMgtWL3QzYbvUjDZTFsLH3et8G1DXM/xA=="
		buf.WriteString(`<script type="application/javascript" src="/assets/js/` + asset.Name + `.js"></script>`)
	}
	return buf.String()
}

// GetStyles returns all assets that should be displayed on a page
func (b *Bundler) GetStyles(page string) string {
	buf := &bytes.Buffer{}
	for _, asset := range b.GetAssets(page) {
		// @TODO: CACHE IN MEMORY
		// @TODO: meta data
		// media="all"
		// crossorigin="anonymous"
		// integrity="sha512-UXiu4O52iBFkqt6Kx5t+pqHYP2/LWWIw9+l5ia74TWw+xPzpH44BFfAQp7yzCe0XFGZa72Xiqyml6tox1KkUjw=="
		buf.WriteString(`<link rel="stylesheet" href="/assets/css/` + asset.Name + `.css">`)
	}
	return buf.String()
}

// build faz o build do bundle, usa DAG para identificar recursos comuns e maximizar a performance de carregamento
// de assets
func (b *Bundler) build() {
	// https://devdocs.magento.com/guides/v2.4/performance-best-practices/advanced-js-bundling.html
	// https://towardsdatascience.com/network-graphs-for-dependency-resolution-5327cffe650f
	// https://ipython-books.github.io/143-resolving-dependencies-in-a-directed-acyclic-graph-with-a-topological-sort/
	// directed acyclic graph (DAG)
	// https://github.com/autom8ter/dagger
	b.assetByName = map[string]*cmn.Asset{}
	b.bundleByPageBuild = map[string][]*cmn.Asset{}
	for page, assets := range b.assetByPage {
		// @TODO: Implementar logica de build correta, por hora só está copiando os assets, é necessário computar dependencias
		b.bundleByPageBuild[page] = assets

		for _, asset := range assets {
			if strings.HasSuffix(asset.Name, ".js") {
				asset.Name = asset.Name[:len(asset.Name)-2]
			} else if strings.HasSuffix(asset.Name, ".css") {
				asset.Name = asset.Name[:len(asset.Name)-3]
			}
			for {
				if byName, exists := b.assetByName[asset.Name]; exists && byName != asset {
					// conflito de nomes de assets, se houver duplicidade, faz a resolução adicionando um sufixo
					// isso não gera problema pois esse processo está sendo realizado em tempo de compilação, até o momento
					// nenhuma página fez mensão a esse arquivo
					asset.Name = asset.Name + "-" + sht.HashXXH64(asset.Name)
				}
				break
			}
			b.assetByName[asset.Name] = asset
		}
	}
}

func (b *Bundler) GetAssetByName(name string) *cmn.Asset {
	return b.assetByName[name]
}

//type AssetType uint8
//
//const (
//	STYLESHEET AssetType = iota
//	JAVASCRIPT
//	IMAGE_PNG
//)

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
