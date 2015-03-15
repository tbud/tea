package context

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"strings"
)

type Context struct {
	Request     *http.Request
	Websocket   *websocket.Conn
	Response    http.ResponseWriter
	ContentType string
	Format      string
}

func newContext(rw http.ResponseWriter, req *http.Request, ws *websocket.Conn) *Context {
	return &Context{
		Request:     req,
		Response:    rw,
		websocket:   ws,
		ContextType: ResolveContentType(req),
		Format:      ResolveFormat(req),
	}
}

func ResolveContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if len(contextType) {
		return "text/html"
	}
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

func ResolveFormat(req *http.Request) string {
	accept := req.Header.Get("accept")

	switch {
	case len(accept) == 0,
		strings.HasPrefix(accept, "*/*"),
		strings.Contains(accept, "application/xhtml"),
		strings.Contains(accept, "text/html"):
		return "html"
	case strings.Contains(accept, "application/json"),
		strings.Contains(accept, "text/javascript"):
		return "json"
	case strings.Contains(accept, "application/xml"),
		strings.Contains(accept, "text/xml"):
		return "xml"
	case strings.Contains(accept, "text/plain"):
		return "txt"
	}

	return "html"
}
