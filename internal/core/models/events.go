package models

type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}
