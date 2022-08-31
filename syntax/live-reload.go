package syntax

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/shtml/sht"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://github.com/cortesi/devd
// https://github.com/sjansen/watchman
// https://github.com/fsnotify/fsnotify
// https://github.com/rollup/rollup/tree/master/src/watch

type liveReloadClient struct {
	addr   string
	events chan *liveReloadEvent
}

type liveReloadEvent struct {
	EventType uint
}

//go:embed static/js/stx-livereload.js
var liveReloadJS embed.FS

// liveReloadInit initialize the site's live-reload client integration
func (s *Site) liveReloadInit(config ConfigLiveReload) error {

	endpoint := strings.TrimSpace(config.Endpoint)
	if endpoint == "" {
		endpoint = "/dev.livereload"
	}

	// add live-reload.js asset, required on all pages
	s.AddFileSystemEmbed(liveReloadJS, "static/", 100)
	asset, err := s.Template.(*sht.TemplateSystem).RegisterAssetJsFilepath("/js/stx-livereload.js")
	if err != nil {
		return err
	}
	s.Bundler.AddRequiredAsset(asset)

	asset.Attributes = map[string]string{
		"data-interval": strconv.Itoa(config.Interval),
		"data-endpoint": endpoint,
	}
	if config.ReloadPageOnCss {
		asset.Attributes["data-reload-page-on-css"] = "true"
	} else {
		asset.Attributes["data-reload-page-on-css"] = "false"
	}

	s.Router.GET(endpoint, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		client := &liveReloadClient{
			addr:   r.RemoteAddr,
			events: make(chan *liveReloadEvent, 10),
		}
		go updateClient(client)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// test
		timeout := time.After(5 * time.Second)

		select {
		case ev := <-client.events:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.Encode(ev)
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	})
	return nil
}

func updateClient(client *liveReloadClient) {
	for {
		client.events <- &liveReloadEvent{
			EventType: uint(rand.Uint32()),
		}
	}
}
