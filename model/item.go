package model

type Item struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"isCompleted"`
}
