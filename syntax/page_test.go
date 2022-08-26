package syntax

import (
	"github.com/syntax-framework/shtml/sht"
	"testing"
)

//
func Test_Page_Directive(t *testing.T) {

	template := `
    <page title="My First Syntax Page - {title}" layout="{layout}"/>
    <div>{valueOne ? 'value-true' : 'value-false'}</div>
  `

	static := []string{
		"",
		"<div>",
		"</div>",
	}

	expected := `<div>value-true</div>`

	values := map[string]interface{}{
		"valueOne": true,
		"title":    "My Dynamic Title",
		"layout":   "custom.html",
	}

	directives := &sht.Directives{}
	directives.Add(PageDirective)

	// expects the directive to extract information at compile time
	compiled, compiler := sht.TestCompile(t, template, static, directives)
	page := compiler.Context.Get(PageConfigKey)
	if page == nil {
		t.Errorf("compiler.Context.Get(PageConfigKey) | invalid output\n   actual: nil expected: *PageConfig")
	} else {
		if compilerPageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
			if compilerPageConfig.Layout != LayoutDefault {
				t.Errorf(
					"compilerPageConfig.Layout | invalid output\n   actual: %s expected: %s", compilerPageConfig.Layout, LayoutDefault,
				)
			}
		} else {
			t.Errorf("compiler.Context.Get(PageConfigKey) | invalid output\n   actual: %v expected: *PageConfig", page)
		}
	}

	_, scope := sht.TestRender(t, compiled, values, expected)

	// expects developer to be able to change page settings at runtime such as title and layout
	page = scope.Context.Get(PageConfigKey)
	if page == nil {
		t.Errorf("scope.Context.Get(PageConfigKey) | invalid output\n   actual: nil expected: *PageConfig")
	} else {
		if renderedPageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
			if renderedPageConfig.Layout != "custom.html" {
				t.Errorf("renderedPageConfig.Layout | invalid output\n   actual: %s expected: custom.html", renderedPageConfig.Layout)
			}
			expectedTitle := "My First Syntax Page - My Dynamic Title"
			if renderedPageConfig.Title != expectedTitle {
				t.Errorf("renderedPageConfig.Title | invalid output\n   actual: %s expected: %s", renderedPageConfig.Title, expectedTitle)
			}
		} else {
			t.Errorf("scope.Context.Get(PageConfigKey) | invalid output\n   actual: %v expected: *PageConfig", page)
		}
	}
}
