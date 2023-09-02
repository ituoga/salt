package salt

import (
	"encoding/json"
	"net/url"
)

type Request struct {
	msg SystemMessage
	ctx *Context
}

func (r *Request) GetPayload() string {
	return r.ctx.req.Payload
}

func (r *Request) Arg() string {
	return r.ctx.req.Arg
}

func (r *Request) To(to interface{}) error {
	return json.Unmarshal([]byte(r.ctx.req.Payload), to)
}

func (r *Request) Query() *url.Values {
	return &r.ctx.req.Query
}
