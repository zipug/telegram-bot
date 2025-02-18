package models

type Answer struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
