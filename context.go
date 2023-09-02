package salt

import (
	"github.com/nats-io/nats.go"
)

// HandlerFunc defines the handler function
type HandlerFunc func(*Context)

type Context struct {
	nc     *nats.Conn
	ncMsg  *nats.Msg
	req    SystemMessage
	rbac   RBACContext
	router *Router
}

func (c *Context) Can(action string) bool {
	if c.rbac == nil {
		return true
	}
	if c.req.Meta == nil {
		return true
	}
	return c.rbac.WithToken(c.req.Meta.Token).WithResource(c.router.resource).Can(action)
}

func (c *Context) Request() *Request {
	return &Request{ctx: c}
}

func (c *Context) Response() *Response {
	return &Response{ctx: c}
}
