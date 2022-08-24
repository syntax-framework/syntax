package syntax

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/shtml"
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
}

// New creates a new site and registers it in mux to handle requests for host.
// If host is the empty string, the registrations are for the wildcard host.
func New(viewsFS fs.FS, viewsBaseDir string) *Site {
	router := httprouter.New()
	site := &Site{
		fsys:         viewsFS,
		viewsBaseDir: viewsBaseDir,
		Router:       router,
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

	// @TODO: Fazer em dois steps, primeiro passeia pelos arquivos, segundo, invoca o parse html ou asset
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

		if strings.HasSuffix(filepath, ".html") {
			errHtmlPage := s.compileHtmlPage(path, templateSystem)
			if errHtmlPage != nil {
				return errHtmlPage // @TODO: Custom error
			}
		} else if isAsset(d.Name()) {
			s.parseAsset(filepath, path, s.fsys)
		}

		return nil
	})
}

// compileHtmlPage load, compile and route page
func (s *Site) compileHtmlPage(path string, ts shtml.TemplateSystem) error {

	pageCompiled, compileContext, err := ts.Compile(path)
	if err != nil {
		return err
	}

	if page := compileContext.Get(PageConfigKey); page != nil {
		if pageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
			// if we have page setup at compile time, you already do layout processing
			println(pageConfig)
		}
	}

	if strings.HasSuffix(path, "/index.html") {
		path = strings.TrimSuffix(path, "index.html")
	}

	s.Router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// @TODO: LastModified, checkPreconditions

		scope := ts.NewScope()

		// compile page content
		pageRendered := pageCompiled.Exec(scope)

		// get page info
		if page := scope.Context.Get(PageConfigKey); page != nil {
			if pageConfig, isPageConfig := page.(*PageConfig); isPageConfig {
				// handles dynamic page parameters, such as title and assets
				println(pageConfig)
			}
		}

		res := pageRendered.String()

		w.Header().Set("GetContent-Restrict", "text/html; charset=utf-8")
		// w.Header().StringSet("GetContent-Length", strconv.FormatInt(sendSize, 10))
		w.WriteHeader(200)

		if r.Method != "HEAD" {
			w.Write([]byte(res))
		}
	})

	return nil
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
