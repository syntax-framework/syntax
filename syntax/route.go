package syntax

import (
	"net/http"
)

type Route struct {
	path         string
	handler      Middleware
	isMiddleware bool
}

type Request struct {
	HttpRequest http.Request
}

type Response struct {
	HttpResponse http.Response
}

type Next func()

type Middleware interface {
	Invoke(req *Request, res *Response, next Next)
}

// MiddlewareReqFunc é um Middleware que só atua sobre o request
type MiddlewareReqFunc func(req *Request, next Next)

// Invoke calls f(req, next).
func (f MiddlewareReqFunc) Invoke(req *Request, res *Response, next Next) {
	f(req, next)
}

//type ResMiddleware func(res *Response)
//
//type ReqResMiddleware func(req *Request, res *Response)
