package salt

import (
	"encoding/json"
)

type Response struct {
	ctx *Context
}

func (r *Response) Reply(payload interface{}, opts ...SystemResponseOption) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	m := &SystemMessage{
		Payload: string(b),
	}
	sr := &SystemResponse{
		Code:    200,
		Message: "OK",
		SysMsg:  m,
	}

	for _, opt := range opts {
		opt(sr)
	}

	b, err = json.Marshal(sr)
	if err != nil {
		return err
	}
	return r.ctx.nc.Publish(r.ctx.ncMsg.Reply, b)
}
