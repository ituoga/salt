package salt

import "encoding/json"

type SystemResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	SysMsg  *SystemMessage `json:"data,omitempty"`
}

func NewSystemResponse(code int, message string) *SystemResponse {
	return &SystemResponse{Code: code, Message: message}
}

func (e *SystemResponse) ToJSON() string {
	// Convert the struct to JSON, handle errors as needed
	json, _ := json.Marshal(e)
	return string(json)
}

// Add this to your salt.Context
func (ctx *Context) ErrorWith(code int, message string) {
	ctx.Response().Reply(NewSystemResponse(code, message).ToJSON())
}

func (ctx *Context) Error(opts ...SystemResponseOption) {
	sr := &SystemResponse{
		Code:    200,
		Message: "OK",
	}
	for _, opt := range opts {
		opt(sr)
	}
	ctx.Response().Reply(sr.ToJSON())
}

type SystemResponseOption func(*SystemResponse)

func WithError(code int, message string) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Code = code
		r.Message = message
	}
}

func WithCode(code int) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Code = code
	}
}

func WithMessage(msg string) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Message = msg
	}
}
