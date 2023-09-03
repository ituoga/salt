package salt

import "encoding/json"

// SystemMessage is the data field of SystemResponse
type SystemResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Type    string         `json:"type,omitempty"`
	SysMsg  *SystemMessage `json:"data,omitempty"`
}

// SystemMessage is the data field of SystemResponse
func NewSystemResponse(code int, message string) *SystemResponse {
	return &SystemResponse{Code: code, Message: message}
}

// Convert the struct to JSON, handle errors as needed
func (e *SystemResponse) ToJSON() string {
	json, _ := json.Marshal(e)
	return string(json)
}

// Add this to your salt.Context
func (ctx *Context) ErrorWith(code int, message string) {
	ctx.Response().Reply(NewSystemResponse(code, message).ToJSON())
}

// Error is a helper function to reply with a SystemResponse and error
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

// SystemResponseOption is a function that modifies a SystemResponse
type SystemResponseOption func(*SystemResponse)

// WithError is a helper function to reply with a SystemResponse and error
func WithError(code int, message string) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Code = code
		r.Message = message
	}
}

// WithCode is a helper function to reply with a SystemResponse and specific code
func WithCode(code int) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Code = code
	}
}

// WithMessage is a helper function to reply with a SystemResponse and specific message
func WithMessage(msg string) SystemResponseOption {
	return func(r *SystemResponse) {
		r.Message = msg
	}
}
