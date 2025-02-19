package models

type Article struct {
	Id          int64
	Name        string
	Description string
	ArticleUrl  string
	Content     string
	ProjectId   int64
	Attachments []Attachment
}
