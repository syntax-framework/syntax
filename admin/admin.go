package admin

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/syntax/syntax"
)

//go:embed views/*
var viewsDir embed.FS

func SiteAdmin() *httprouter.Router {
	adm := syntax.New(viewsDir, "views")

	adm.Init()

	return adm.Router
}
