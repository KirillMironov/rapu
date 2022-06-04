package domain

import (
	"encoding/json"
)

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
