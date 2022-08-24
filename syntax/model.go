package syntax

import (
	"context"
	"github.com/julienschmidt/httprouter"
)

//type ModelResult struct {
//  Data  interface{}
//  Cache interface{}
//}

type ModelResult struct {
	Data  interface{}
	Cache interface{}
}

func (r ModelResult) WithCache() {

}

func (r ModelResult) CaseNotInCache(func()) {

}

type Model interface {
	Before(request Request, ctx context.Context, params httprouter.Params) ModelResult
	//Data(request Request, ctx context.Context, params httprouter.Params, ModelResult) ModelResult
	//Procces(request Request, ctx context.Context, params httprouter.Params, Config) ModelResult
}

//func Prepare(params) ModelResult {
//  return ModelResult{
//    Cache: {
//      id: "pxoto" + pamras.id,
//    },
//  }
//}
