package context

type Router struct {
	Routes []*Router
	path   string
}

func (r *Router) Load() {

}
