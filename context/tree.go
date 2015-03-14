package router

type Handle func() error

type nodeType uint8

const (
	static_type nodeType = iota
	param_type
	match_type
	start_type
)

type node struct {
	nType    nodeType
	hasChild bool
	children []*node
	handle   Handle
}
