package syntax

import (
	"github.com/syntax-framework/shtml/sht"
)

const PageConfigKey = "syntax.page.config"

// PageConfig configuration of a page in syntax framework
type PageConfig struct {
	Layout string // name of the layout file used to render this page
	Title  string // page title
}

// PageDirective extrai as configurações da página a partir da tag <page />
var PageDirective = &sht.Directive{
	Name:       "page",
	Restrict:   sht.ELEMENT,
	Priority:   100,
	Terminal:   true,
	Transclude: nil,
	Compile: func(node *sht.Node, attrs *sht.Attributes, t *sht.Compiler) (*sht.DirectiveMethods, error) {
		return &sht.DirectiveMethods{
			Process: func(scope *sht.Scope, attrs *sht.Attributes, _ sht.TranscludeFunc) *sht.Rendered {
				pageConfig := &PageConfig{
					Layout: attrs.GetOrDefault("layout", "root"),
					Title:  attrs.Get("title"),
				}
				scope.Context.Set(PageConfigKey, pageConfig)
				return nil
			},
		}, nil
	},
}
