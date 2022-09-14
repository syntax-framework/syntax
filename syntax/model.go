package syntax

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
	Before() ModelResult
	//Data(request Request, ctx context.Context, ModelResult) ModelResult
	//Procces(request Request, ctx context.Context, Config) ModelResult
}

//func Prepare(params) ModelResult {
//  return ModelResult{
//    Cache: {
//      id: "pxoto" + pamras.id,
//    },
//  }
//}
