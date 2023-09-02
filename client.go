package salt

import (
	"encoding/json"
	"io"
	"net/url"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type Client struct {
	nc *nats.Conn
}

func NewClient(natsURL string) (*Client, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &Client{nc: nc}, nil
}

func (c *Client) Close() {
	c.nc.Close()
}

type ResponseToClient struct {
	Data    responseToClient `json:"data,omitempty"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
	// SysMsg  *SystemMessage `json:"data,omitempty"`
}

type responseToClient struct {
	Payload string `json:"payload"`
}

func (r *ResponseToClient) CopyTo(to *json.RawMessage) error {
	rd := r.Data
	*to = []byte(rd.Payload)
	return nil
}

func (r *ResponseToClient) Write(w io.Writer) (int, error) {
	return io.WriteString(w, r.Data.Payload)
}

func (r *ResponseToClient) To(to interface{}) error {
	return json.Unmarshal([]byte(r.Data.Payload), to)
}

func (r *ResponseToClient) from(data []byte) error {
	return json.Unmarshal(data, r)
}

type ClientRequestOpts func(*SystemMessage)

func WithPayloadClientRaw(payload interface{}) ClientRequestOpts {
	return func(sm *SystemMessage) {
		sm.Payload = string(payload.([]byte))
	}
}

func WithPayloadClient(payload interface{}) ClientRequestOpts {
	return func(sm *SystemMessage) {
		b, _ := json.Marshal(payload)
		sm.Payload = string(b)
	}
}

func WithTokenClient(token string) ClientRequestOpts {
	return func(sm *SystemMessage) {
		sm.Meta = &Meta{
			Token: token,
		}
	}
}

func WithQueryClient(uv url.Values) ClientRequestOpts {
	return func(sm *SystemMessage) {
		sm.Query = uv
	}
}

func WithArg(args string) ClientRequestOpts {
	return func(sm *SystemMessage) {
		sm.Arg = args
	}
}

func (c *Client) Request(topic string, opts ...ClientRequestOpts) (*ResponseToClient, error) {
	req := SystemMessage{}
	for _, opt := range opts {
		opt(&req)
	}
	b, err := req.Marshal()
	if err != nil {
		return nil, errors.Wrapf(err, "Client request -> marshalling request: %v", req)
	}
	msg, err := c.nc.Request(topic, b, 5*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "Client request -> sending request %v", topic)
	}
	res := ResponseToClient{}
	err = res.from(msg.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "Client request -> unmarshalling response: %v", msg.Data)
	}
	return &res, nil
}
