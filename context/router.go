package router

import (
	"regexp"
)

type Router struct {
	Routes []*Route
	path   string
}

func (r *Router) Load() {

}
