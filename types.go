package main

type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeEmail ContentType = "email"
)

type Message struct {
	Content 	string 		`json:"content"`
	ContentType ContentType `json:"content_type"`
}

type SpamResult struct {
	IsSpam bool `json:"is_spam"`
	Score int `json:"score"`
	Reason string `json:"reason"`
	ContentType ContentType `json:"content_type"`
}