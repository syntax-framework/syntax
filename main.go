package main

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/syntax/syntax"
	"log"
	"net/http"
	"os"
)

// @TODO: https://github.com/fsnotify/fsnotify
func main() {
	handler := createSite()
	httpAddr := "localhost:8080"
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", httpAddr, err)
	}
}

//go:embed site_embed/*
var embedSiteDir embed.FS

func createSite() *httprouter.Router {

	config := &syntax.Config{
		Dev: true,
		LiveReload: syntax.ConfigLiveReload{
			Interval: 100,
			Debounce: 200,
			//ReloadPageOnCss: false,
			//Patterns:        nil,
			//Endpoint:        "",
		},
	}

	site := syntax.New(config)

	if config.Dev {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		site.AddFileSystemDir(path+"/site_live_reload/", 0)
	}

	//site.AddFileSystemEmbed(embedSiteDir, "site_embed/", 0) // test only
	//site.AddFileSystemEmbed(embedSiteDir2, "site_test/pages/", 0) // test only
	//site.AddFileSystemEmbed(embedSiteDir3, "site_test/pages/", 0) // test only

	if err := site.Init(); err != nil {
		log.Fatal(err)
	}

	return site.Router
}
