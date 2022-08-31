package syntax

import (
	"bytes"
	"github.com/syntax-framework/shtml/cmn"
	"github.com/syntax-framework/shtml/sht"
)

// Bundler responsible for grouping the assets used by the pages and, from that, defining the most optimized way to
// group these resources in order to maximize performance.
type Bundler struct {
	dirty             bool                    // Indica que houve mudança na
	assetByName       map[string]*cmn.Asset   // facilita busca
	assetByPage       map[string][]*cmn.Asset // Lista original de assets por página
	assetRequired     map[*cmn.Asset]bool     // facilita busca
	bundleByPageBuild map[string][]*cmn.Asset // Lista processada de assets por página
}

func (b *Bundler) AddRequiredAsset(asset *cmn.Asset) {
	if b.assetRequired == nil {
		b.assetRequired = map[*cmn.Asset]bool{}
	}
	b.assetRequired[asset] = true
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
func (b *Bundler) GetAssets(page string, assetType cmn.AssetType) []*cmn.Asset {
	if b.dirty {
		b.build()
		b.dirty = false
	}

	var assets cmn.Assets
	if bundle, exists := b.bundleByPageBuild[page]; exists {
		for _, asset := range bundle {
			if asset.Type == assetType {
				assets = append(assets, asset)
			}
		}
	}

	// required assets
	for asset, _ := range b.assetRequired {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}

	dependencies, err := assets.Resolve()
	if err != nil {
		// @TODO: Ciclic dependencies, how to solve?
		println(err)
		return nil
	}

	return dependencies
}

// GetScripts returns all assets that should be displayed on a page
func (b *Bundler) GetScripts(page string) string {
	buf := &bytes.Buffer{}
	for _, asset := range b.GetAssets(page, cmn.Javascript) {
		// @TODO: meta data
		// defer="defer"
		// crossorigin="anonymous"
		buf.WriteString(`<script type="application/javascript"`)

		if asset.Url != "" {
			buf.WriteString(` src="` + asset.Url + `"`)
		} else {
			buf.WriteString(` src="/assets/js/` + asset.Name + `.js"`)
		}

		if asset.Integrity != "" {
			buf.WriteString(` integrity="` + asset.Integrity + `"`)
		}

		if asset.Attributes != nil {
			for name, value := range asset.Attributes {
				buf.WriteString(name)
				if value != "" {
					buf.WriteString(`="` + sht.HtmlEscape(value) + `"`)
				}
			}
		}

		buf.WriteString(`></script>`)
	}
	return buf.String()
}

// GetStyles returns all assets that should be displayed on a page
func (b *Bundler) GetStyles(page string) string {
	buf := &bytes.Buffer{}
	for _, asset := range b.GetAssets(page, cmn.Stylesheet) {
		// @TODO: CACHE IN MEMORY
		// @TODO: meta data
		// media="all"
		// crossorigin="anonymous"
		buf.WriteString(`<link rel="stylesheet"`)

		if asset.Url != "" {
			buf.WriteString(` href="` + asset.Url + `"`)
		} else {
			buf.WriteString(` href="/assets/css/` + asset.Name + `.css"`)
		}

		if asset.Integrity != "" {
			buf.WriteString(` integrity="` + asset.Integrity + `"`)
		}

		if asset.Attributes != nil {
			for name, value := range asset.Attributes {
				buf.WriteString(" " + name)
				if value != "" {
					buf.WriteString(`="` + sht.HtmlEscape(value) + `"`)
				}
			}
		}

		buf.WriteString(`>`)
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
			b.assetByName[asset.Name] = asset
		}
	}

	for asset, _ := range b.assetRequired {
		b.assetByName[asset.Name] = asset
	}
}

func (b *Bundler) GetAssetByName(name string) *cmn.Asset {
	return b.assetByName[name]
}
