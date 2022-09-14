package syntax

import (
	"github.com/iancoleman/strcase"
	"github.com/syntax-framework/shtml/cmn"
	"github.com/syntax-framework/shtml/sht"
	"strings"
)

var errorControllerNotFound = cmn.Err(
	"controller.notfound",
	"There is no controller registered with the given name.", "Name: %s", "Component: %s",
)

func (s *Syntax) CreateControllerDirectives() []*sht.Directive {

	elementDirective := &sht.Directive{
		Name:       "controller",
		Restrict:   sht.ATTRIBUTE,
		Priority:   200,
		Terminal:   false,
		Transclude: false,
		Compile: func(node *sht.Node, attrs *sht.Attributes, t *sht.Compiler) (methods *sht.DirectiveMethods, err error) {

			name := attrs.Get("controller")
			attrs.Remove(attrs.GetAttribute("controller"))

			var controller *Controller
			for _, ctrl := range s.Controllers {
				if ctrl.Name == name {
					controller = ctrl
				}
			}

			if controller == nil {
				err = errorControllerNotFound(name, node.DebugTag())
				return
			}

			methods = &sht.DirectiveMethods{
				Process: func(scope *sht.Scope, attrs *sht.Attributes, transclude sht.TranscludeFunc) *sht.Rendered {
					params := map[string]interface{}{}
					for attrName, attr := range attrs.Map {
						if strings.HasPrefix(attrName, "param-") {
							paramName := strcase.ToLowerCamel(strings.Replace(attrName, "param-", "", 1))
							params[paramName] = attr.Value
							attrs.Remove(attr)
						}
					}

					controller.Setup(scope, params)

					if controller.Live != nil {
						// is live controller
						// serialize params to allow reconnection
						attrs.Set("data-stx-ctrl", controller.Name)
						attrs.Set("data-stx-ctrl-par", "parametros serializado, todo")
					}

					return nil
					//return transclude("", nil)
				},
			}

			return
		},
	}

	// ControllerDirective extrai as configurações da página a partir da tag <page />
	return []*sht.Directive{elementDirective}
}
