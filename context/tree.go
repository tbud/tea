package context

import (
	"net/http"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type nodeType uint8

type Param struct {
	Key   string
	Value string
}

type Params []Param

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
