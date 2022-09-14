package syntax

import (
	"github.com/syntax-framework/shtml/cmn"
	"github.com/syntax-framework/shtml/sht"
	"log"
	"strings"
)

var errorControllerInvalidState = cmn.Err(
	"controller.invalid.state",
	"It is not allowed to register new controllers after initialization.", "Name: %s",
)

// se desenvolvedor precisar de dados da requisição, criar um midleware e adiconar no Context.
// O contexto fica disponível em toda as arvore de execução
//type HttpCtx struct {
//	Request *http.Request
//}

type Controller struct {
	Name  string
	Setup ControllerSetupFunc
	Live  ControllerLiveFunc
}

type LiveState struct {
}

// On p
func (l LiveState) On(event string, callback func(params map[string]interface{})) {

}

type ControllerSetupFunc func(scope *sht.Scope, params map[string]interface{})

type ControllerLiveFunc func(scope *sht.Scope, params map[string]interface{}, live *LiveState)

func (s *Syntax) RegisterController(name string, setup ControllerSetupFunc, live ControllerLiveFunc) {

	if s.initialized {
		failToStart(
			errorControllerInvalidState(name).Error(),
			"Change the controller record to run before the method `syntax.Init()`",
		)
	}

	name = strings.TrimSpace(name)
	if name == "" || strings.ContainsRune(name, ' ') {
		log.Fatal("Nome de controller inválido " + name)
	}

	for _, controller := range s.Controllers {
		if controller.Name == name {
			log.Fatal("Já existe uma controller com o nome " + name)
		}
	}

	controller := &Controller{
		Name:  name,
		Setup: setup,
		Live:  live,
	}
	s.Controllers = append(s.Controllers, controller)
}
