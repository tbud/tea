package context

import (
	"net/http"
)

type Context struct {
	Request  *Request
	Response *Response
}

type Request struct {
	*http.Request
	ContextType string
	Format      string
}

type Response struct {
	Status      int
	ContentType string

	Out http.ResponseWriter
}
