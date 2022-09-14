package syntax

import (
	"github.com/syntax-framework/shtml/sht"
	"strings"
)

const PageConfigKey = "syntax.page.config"
const LayoutDefault = "root"

// PageConfig configuration of a page in syntax framework
type PageConfig struct {
	Layout string // name of the layout file used to render this page
	Title  string // page title
}

func checkPageConfig(config *PageConfig) {
	if strings.ContainsAny(config.Layout, "{") {
		config.Layout = LayoutDefault
	}

	if strings.ContainsAny(config.Title, "{") {
		config.Title = ""
	}
}

// PageDirective extrai as configurações da página a partir da tag <page />
var PageDirective = &sht.Directive{
	Name:       "page",
	Restrict:   sht.ELEMENT,
	Priority:   200,
	Terminal:   true,
	Transclude: true,
	Compile: func(node *sht.Node, attrs *sht.Attributes, t *sht.Compiler) (*sht.DirectiveMethods, error) {

		compileConfig := &PageConfig{
			Layout: layoutValidName(attrs.GetOrDefault("layout", LayoutDefault)),
			Title:  attrs.Get("title"),
		}
		checkPageConfig(compileConfig)
		t.Context.Set(PageConfigKey, compileConfig)

		return &sht.DirectiveMethods{
			Process: func(scope *sht.Scope, attrs *sht.Attributes, _ sht.TranscludeFunc) *sht.Rendered {
				runtimeConfig := &PageConfig{
					// Dynamic part
					Title: attrs.Get("title"),
					// Static config (compile time)
					Layout: compileConfig.Layout,
				}
				checkPageConfig(runtimeConfig)
				scope.Context.Set(PageConfigKey, runtimeConfig)
				return nil
			},
		}, nil
	},
}
