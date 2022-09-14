package syntax

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/syntax-framework/shtml/sht"
	"io"
	"net/http"
	"time"
)

var fallbackSessionIDSeq = &sht.Sequence{Salt: time.Now().String()}

// RegenerateSessionID permite gerar um novo session ID para o usuario
//
// https://owasp.org/www-community/attacks/Session_fixation
func (s *Syntax) RegenerateSessionID(w http.ResponseWriter, r *http.Request) {

	var sid string
	bytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		sid = sht.HashXXH64Hex(fallbackSessionIDSeq.NextHash() + time.Now().String())
	} else {
		sid = base64.URLEncoding.EncodeToString(bytes)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.Config.Cookie.Name,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   s.Config.Cookie.MaxAge,
	})
}

func (s *Syntax) Use(args ...interface{}) {
	s.router.Use(args...)
}

func (s *Syntax) GET(path string, handle interface{}) {
	s.router.GET(path, handle)
}

func (s *Syntax) HEAD(path string, handle interface{}) {
	s.router.HEAD(path, handle)
}

func (s *Syntax) OPTIONS(path string, handle interface{}) {
	s.router.OPTIONS(path, handle)
}

func (s *Syntax) POST(path string, handle interface{}) {
	s.router.POST(path, handle)
}

func (s *Syntax) PUT(path string, handle interface{}) {
	s.router.PUT(path, handle)
}

func (s *Syntax) PATCH(path string, handle interface{}) {
	s.router.PATCH(path, handle)
}

func (s *Syntax) DELETE(path string, handle interface{}) {
	s.router.DELETE(path, handle)
}

func (s *Syntax) Handle(method string, path string, handle interface{}) {
	s.router.Handle(method, path, handle)
}
