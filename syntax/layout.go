package syntax

import (
	"github.com/syntax-framework/shtml"
	"github.com/syntax-framework/shtml/sht"
	"strings"
)

/**
Métodos para gerenciamento de layouts de uma página
*/

// Layout referencias de um layout compilado
type Layout struct {
	Name     string
	Compiled *sht.Compiled
}

// getLayout obtém um layout por nome
func getLayout(name string, ts shtml.TemplateSystem) (*Layout, error) {
	compiled, _, err := ts.Compile("/_layout/" + name)
	if err != nil {
		return nil, err
	}

	return &Layout{
		Name:     name,
		Compiled: compiled,
	}, nil
}

func layoutValidName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return LayoutDefault
	}
	if !strings.HasSuffix(name, ".html") {
		return name + ".html"
	}
	return name
}
