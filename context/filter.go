package context

type Filter func(c *Context, filterChain []Filter)
