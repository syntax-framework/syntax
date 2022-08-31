package syntax

import (
	"bytes"
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/syntax-framework/shtml"
	"github.com/syntax-framework/shtml/cmn"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

type RouteType uint8

const (
	ASSET RouteType = 1 << iota
	PAGE
	MODEL
)

// FileSystem referencia para um fs.FS, permite que o sistema trabalhe com vários FileSystem
type FileSystem struct {
	fsys     fs.FS
	root     string   // root dir
	Ignored  []string // existing directories in that fsys starting with `_`
	priority int      // allows prioritizing filesystems, used by libs to make components available
	embed    bool
}

type Site struct {
	Router      *httprouter.Router
	Config      *Config
	Bundler     *Bundler
	FileSystems []*FileSystem
	//host       string
	pages        []*PageConfig
	models       []*Model
	middleware   []*Middleware
	viewsBaseDir string
	filesLookup  map[string]*FileSystem // cache lookup
	Template     shtml.TemplateSystem
}

// New creates a new site and registers it in mux to handle requests for host.
// If host is the empty string, the registrations are for the wildcard host.
func New(config *Config) *Site {
	router := httprouter.New()
	site := &Site{
		//fsys:         viewsFS,
		//viewsBaseDir: viewsBaseDir,
		Config:  config,
		Router:  router,
		Bundler: &Bundler{},
		//host:   host,
		filesLookup: map[string]*FileSystem{},
	}

	site.Template = shtml.New(func(filepath string) (string, error) {
		return site.loadFile(filepath)
	})

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

	s.registerDirectives()

	config := s.Config
	if config.Dev {
		// live reload
		if !config.LiveReload.Disabled {
			s.liveReloadInit(config.LiveReload)
		}
	}

	s.serveAssets()

	// serve pages
	err := s.servePages()
	if err != nil {
		return err
	}

	return nil
}

// AddFileSystemDir register a new directory FileSystem on that site
func (s *Site) AddFileSystemDir(root string, priority int) {
	dir := path.Clean(root)
	s.addFileSystem(&FileSystem{
		fsys:     os.DirFS(dir),
		priority: priority,
		root:     "",
		embed:    false,
	})
}

// AddFileSystemEmbed register a new embed FileSystem on that site
func (s *Site) AddFileSystemEmbed(embedFs embed.FS, root string, priority int) {
	s.addFileSystem(&FileSystem{
		fsys:     embedFs,
		priority: priority,
		embed:    true,
		root:     path.Clean(root),
	})
}

// AddFileSystemEmbed register a new FileSystem on that site
func (s *Site) addFileSystem(system *FileSystem) {
	s.FileSystems = append(s.FileSystems, system)
	sort.Slice(s.FileSystems, func(i, j int) bool {
		return s.FileSystems[i].priority > s.FileSystems[j].priority
	})
}

// registerDirectives register custom Syntax directives
func (s *Site) registerDirectives() {
	s.Template.Register(PageDirective)
}

func (s *Site) servePages() error {
	for _, system := range s.FileSystems {
		err := fs.WalkDir(system.fsys, ".", func(filepath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if filepath == "." {
				return nil
			}

			if d.IsDir() {
				if strings.HasPrefix(d.Name(), "_") {
					system.Ignored = append(system.Ignored, filepath)
				}
				return nil
			}

			// ignored directories are only used in imports
			for _, ignored := range system.Ignored {
				if strings.HasPrefix(filepath, ignored) {
					return nil
				}
			}

			// only html pages
			if strings.HasSuffix(filepath, ".html") {
				fPath := strings.TrimPrefix(filepath, system.root)
				//if fPath[0] != '/' {
				//	fPath = "/" + fPath
				//}

				errHtmlPage := s.processPage(fPath)
				if errHtmlPage != nil {
					return errHtmlPage // @TODO: Custom error
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// processPage load, compile and route page
func (s *Site) processPage(path string) error {

	pageCompiled, compileContext, err := s.Template.Compile(path)
	if err != nil {
		return err
	}

	if path[0] != '/' {
		path = "/" + path
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
	if layout, err = s.getLayout(layoutName); err != nil {
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
	s.Bundler.SetPageAssets(path, assets)

	s.Router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// @TODO: LastModified, checkPreconditions

		pageConfigRuntime := pageConfigCompile

		// compile page content
		pageScope := s.Template.NewScope()
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

		layoutScope := s.Template.NewScope()
		layoutScope.Set("page", pageConfigRuntime)
		layoutScope.Set("content", pageRendered.String())
		layoutScope.Set("styles", s.Bundler.GetStyles(path))
		layoutScope.Set("scripts", s.Bundler.GetScripts(path))
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
			asset = s.Bundler.GetAssetByName(strings.TrimPrefix(strings.TrimSuffix(filepath, ".css"), "/css/"))
			if asset != nil && asset.Type != cmn.Stylesheet {
				asset = nil
			}
		} else if strings.HasPrefix(filepath, "/js/") {
			// se o asset estiver em um bundle e, o tamanho do arquivo com relação ao bundler for muito menor
			// entregar o conteúdo js, caso contrário, fazer redirecionamento para o bundler
			asset = s.Bundler.GetAssetByName(strings.TrimPrefix(strings.TrimSuffix(filepath, ".js"), "/js/"))
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

// loadFile load a file from FileSystems
func (s *Site) loadFile(filepath string) (string, error) {

	var file fs.File
	var fileSystem *FileSystem
	var err error

	// lookup
	if system, found := s.filesLookup[filepath]; found {
		file, err = system.fsys.Open(system.root + filepath)
		if err != nil {
			if pathError, isPathError := err.(*fs.PathError); isPathError && pathError.Err == fs.ErrNotExist {
				// file removed from this file system
				delete(s.filesLookup, filepath)
			} else {
				return "", err
			}
		}
		fileSystem = system
	}

	if file == nil {
		for _, system := range s.FileSystems {
			fullPath := strings.TrimPrefix(path.Join(system.root, filepath), "/")
			file, err = system.fsys.Open(fullPath)
			if err != nil {
				if pathError, isPathError := err.(*fs.PathError); isPathError && pathError.Err == fs.ErrNotExist {
					continue
				}
				return "", err
			}
			fileSystem = system
			break
		}
	}

	if file == nil {
		// not found
		return "", fs.ErrNotExist
	}

	s.filesLookup[filepath] = fileSystem

	// load content
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
