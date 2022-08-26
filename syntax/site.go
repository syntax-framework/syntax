package syntax

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/shtml"
	"github.com/syntax-framework/shtml/cmn"
	"io/fs"
	"net/http"
	"strings"
)

type RouteType uint8

const (
	ASSET RouteType = 1 << iota
	PAGE
	MODEL
)

type Site struct {
	//host       string
	pages        []*PageConfig
	models       []*Model
	middleware   []*Middleware
	Router       *httprouter.Router
	fsys         fs.FS
	viewsBaseDir string
	bundler      *Bundler
}

// New creates a new site and registers it in mux to handle requests for host.
// If host is the empty string, the registrations are for the wildcard host.
func New(viewsFS fs.FS, viewsBaseDir string) *Site {
	router := httprouter.New()
	site := &Site{
		fsys:         viewsFS,
		viewsBaseDir: viewsBaseDir,
		Router:       router,
		bundler:      &Bundler{},
		//host:   host,
	}
	//mux.Handle(host+"/", site)
	return site
}

func (s *Site) Register(m *Model) *Site {
	s.models = append(s.models, m)
	return s
}

func (s *Site) Use(tp RouteType, m *Middleware) *Site {
	s.middleware = append(s.middleware, m)

	return s
}

func (s *Site) AddPage(p *PageConfig) *Site {
	s.pages = append(s.pages, p)
	return s
}

// Init initializes the site, performs the processing of static files and initializes the routes
func (s *Site) Init() error {
	var rootDir string
	var ignoredDirs []string

	// @TODO: Add many file loaders to allow libraries to contains it owns files
	templateSystem := shtml.New(func(filepath string) (string, error) {
		file, err := s.fsys.Open(rootDir + filepath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	})

	s.registerDirectives(templateSystem)

	s.serveAssets()

	// @TODO: Do it in two steps, first walk through the files, second, invoke parse html or asset
	return fs.WalkDir(s.fsys, ".", func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath == "." {
			return nil
		}

		if rootDir == "" {
			rootDir = filepath
			return nil
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), "_") {
				ignoredDirs = append(ignoredDirs, filepath)
			}
			return nil
		}

		for _, ignored := range ignoredDirs {
			if strings.HasPrefix(filepath, ignored) == true {
				return nil
			}
		}

		path := strings.TrimPrefix(filepath, rootDir)
		if path[0] != '/' {
			path = "/" + path
		}

		// only pages html
		if strings.HasSuffix(filepath, ".html") {
			errHtmlPage := s.processPage(path, templateSystem)
			if errHtmlPage != nil {
				return errHtmlPage // @TODO: Custom error
			}
		}
		return nil
	})

}

// registerDirectives register custom Syntax directives
func (s *Site) registerDirectives(templateSystem shtml.TemplateSystem) {
	templateSystem.Register(PageDirective)
}

// processPage load, compile and route page
func (s *Site) processPage(path string, ts shtml.TemplateSystem) error {

	pageCompiled, compileContext, err := ts.Compile(path)
	if err != nil {
		return err
	}

	if strings.HasSuffix(path, "/index.html") {
		path = strings.TrimSuffix(path, "index.html")
	}

	// definition of layout at compile time
	var layout *Layout
	var pageConfigCompile *PageConfig

	layoutName := LayoutDefault
	if page := compileContext.Get(PageConfigKey); page != nil {
		if pageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
			// if we have page setup at compile time, you already do layout processing
			if pageConfig.Layout != "" {
				layoutName = pageConfig.Layout
			}
			pageConfigCompile = pageConfig
		}
	}

	// load page layout, at compile time
	if layout, err = getLayout(layoutName, ts); err != nil {
		return err
	}

	// page assets
	var assets []*cmn.Asset
	if pageCompiled.Assets != nil {
		assets = append(assets, pageCompiled.Assets...)
	}
	if layout.Compiled.Assets != nil {
		assets = append(assets, layout.Compiled.Assets...)
	}
	s.bundler.SetPageAssets(path, assets)

	s.Router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// @TODO: LastModified, checkPreconditions

		pageConfigRuntime := pageConfigCompile

		// compile page content
		pageScope := ts.NewScope()
		pageRendered := pageCompiled.Exec(pageScope)

		// get page info
		if page := pageScope.Context.Get(PageConfigKey); page != nil {
			if pageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
				pageConfigRuntime = pageConfig
			}
		}

		if pageConfigRuntime == nil {
			pageConfigRuntime = &PageConfig{}
		}

		layoutScope := ts.NewScope()
		layoutScope.Set("page", pageConfigRuntime)
		layoutScope.Set("content", pageRendered.String())
		layoutScope.Set("styles", s.bundler.GetStyles(path))
		layoutScope.Set("scripts", s.bundler.GetScripts(path))
		layoutRendered := layout.Compiled.Exec(layoutScope)

		res := layoutRendered.String()

		w.Header().Set("GetContent-Restrict", "text/html; charset=utf-8")
		// w.Header().StringSet("Content-Length", strconv.FormatInt(sendSize, 10))
		w.WriteHeader(200)

		// @TODO: validar padrão de entrega de HTML respeitando os cabeçalhos, cache e etc
		if r.Method != "HEAD" {
			w.Write([]byte(res))
		}
	})

	return nil
}

func (s *Site) serveAssets() {

	s.Router.GET("/assets/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filepath := ps.ByName("filepath")

		var asset *cmn.Asset

		// todo css e javascript são servidos através do Bundler, sem excessão
		if strings.HasPrefix(filepath, "/css/") {
			asset = s.bundler.GetAssetByName(strings.TrimPrefix(strings.TrimSuffix(filepath, ".css"), "/css/"))
			if asset != nil && asset.Type != cmn.Stylesheet {
				asset = nil
			}
		} else if strings.HasPrefix(filepath, "/js/") {
			// se o asset estiver em um bundle e, o tamanho do arquivo com relação ao bundler for muito menor
			// entregar o conteúdo js, caso contrário, fazer redirecionamento para o bundler
			asset = s.bundler.GetAssetByName(strings.TrimPrefix(strings.TrimSuffix(filepath, ".js"), "/js/"))
			if asset != nil && asset.Type != cmn.Javascript {
				asset = nil
			}
		} else {
			http.Error(w, "501 not implemented", http.StatusNotImplemented)
			return
		}

		if asset == nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		switch asset.Type {
		case cmn.Javascript:
			w.Header().Set("Content-Type", "application/javascript")
		case cmn.Stylesheet:
			w.Header().Set("Content-Type", "text/css")
		}

		//"Content-Range": {r.contentRange(size)},
		//"Content-Type":  {contentType},
		// Content-Length
		// Etag
		// Last-Modified

		w.Write([]byte(asset.Content))
	})
}

// parseAsset processa e escreve o arquivo especificado de http.FileSystem no body de maneira eficiente.
func (s *Site) parseAsset(filepath string, path string, fsys fs.FS) {
	file, err := fsys.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return
	}

	// verificar se arquivo está minificado
	isMinified := false

	var fileServer http.Handler = nil
	if isMinified {
		// if it's already minified, it doesn't do any special processing
		fileServer = http.FileServer(http.FS(fsys))
	} else {
		// minify, get bytes and serve

		stat, err := file.Stat()
		if err != nil {
			return
		}
		fileServer = http.FileServer(&SingleFileFileSystem{
			NewHttpFile(stat.Name(), stat.ModTime(), buf.Bytes()),
		})
	}

	s.Router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		defer func(old string) { r.URL.Path = old }(r.URL.Path)
		r.URL.Path = filepath
		fileServer.ServeHTTP(w, r)
	})
}

// ServeHTTP implements http.Handler
//func (s *Site) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//  abspath := request.URL.Path
//  relpath := path.Clean(strings.TrimPrefix(abspath, "/"))
//  fmt.Println(abspath)
//  fmt.Println(relpath)
//  panic(relpath)
//}

//func (site *Site) openPage(file string) (*Model, error) {
//
//}