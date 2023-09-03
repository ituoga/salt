package salt

import (
	"encoding/json"
)

// Response struct
type Response struct {
	ctx *Context
}

// Notify sends a notification to a topic with payload and options
func (r *Response) Notify(topic string, payload interface{}, opts ...SystemResponseOption) error {
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
	tmpTopic := topic
	if r.ctx.router.resource != "" {
		tmpTopic = r.ctx.router.resource + "." + topic
	}
	return r.ctx.nc.Publish(tmpTopic, b)
}

// Reply sends a reply to the request with payload and options
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
