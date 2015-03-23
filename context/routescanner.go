package context

import (
	"github.com/tbud/x/container/set"
)

type route struct {
	imports map[string]*set.StringSet // save import path with prefix key
	routers []router
}

type router struct {
	httpMethod string
	path       string
	prefix     string
	structName string
	action     string
	params     []param
}

type param struct {
	name         string
	pType        paramType
	defaultValue interface{}
}

type paramType uint8

const (
	default_type paramType = iota
	fix_value_type
	default_value_type
)
