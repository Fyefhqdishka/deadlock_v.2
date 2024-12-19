package models

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
