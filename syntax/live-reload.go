package syntax

import (
	"github.com/syntax-framework/shtml/sht"
	"strconv"
	"strings"
)

// https://github.com/cortesi/devd
// https://github.com/sjansen/watchman
// https://github.com/fsnotify/fsnotify
// https://github.com/rollup/rollup/tree/master/src/watch
// https://github.com/fsnotify/fsnotify

type liveReloadClient struct {
	addr   string
	events chan *liveReloadEvent
}

type liveReloadEvent struct {
	EventType uint
}

// liveReloadInit initialize the site's live-reload client integration
func (s *Syntax) liveReloadInit(config ConfigLiveReload) error {

	endpoint := strings.TrimSpace(config.Endpoint)
	if endpoint == "" {
		endpoint = "/dev.livereload"
	}

	// add live-reload.js asset, required on all pages
	asset, err := s.Template.(*sht.TemplateSystem).RegisterAssetJsFilepath("/assets/js/stx-livereload.js")
	if err != nil {
		return err
	}
	s.Bundler.AddRequiredAsset(asset)

	asset.Attributes = map[string]string{
		"data-interval": strconv.Itoa(config.Interval),
		"data-endpoint": endpoint,
	}
	if config.ReloadCss {
		asset.Attributes["data-reload-page-on-css"] = "true"
	} else {
		asset.Attributes["data-reload-page-on-css"] = "false"
	}

	//s.Router.GET(endpoint, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//	// We need to be able to flush for SSE
	//	flusher, ok := w.(http.Flusher)
	//	if !ok {
	//		http.Error(w, "Connection does not support streaming", http.StatusBadRequest)
	//		return
	//	}
	//
	//	client := &liveReloadClient{
	//		addr:   r.RemoteAddr,
	//		events: make(chan *liveReloadEvent, 10),
	//	}
	//	go updateClient(client)
	//
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//	w.Header().Set("Content-Type", "text/event-stream")
	//	w.Header().Set("Cache-Control", "no-cache")
	//	w.Header().Set("Connection", "keep-alive")
	//
	//	// test
	//	timeout := time.After(5 * time.Second)
	//
	//	select {
	//	case ev := <-client.events:
	//		var buf bytes.Buffer
	//		enc := json.NewEncoder(&buf)
	//		enc.Encode(ev)
	//		fmt.Fprintf(w, "data: %v\n\n", buf.String())
	//		flusher.Flush()
	//	case <-timeout:
	//		fmt.Fprintf(w, ": nothing to sent\n\n")
	//	}
	//})
	return nil
}
