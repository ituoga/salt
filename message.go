package salt

import (
	"encoding/json"
	"net/url"
)

type Meta struct {
	Token string `json:"token"`
}

type SystemMessage struct {
	Meta    *Meta      `json:"meta"`
	Payload string     `json:"payload"`
	Query   url.Values `json:"query"`
	Arg     string     `json:"arg"`
}

func (m *SystemMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *SystemMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
