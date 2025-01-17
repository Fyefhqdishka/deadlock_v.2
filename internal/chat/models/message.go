package models

type Message struct {
	Sender   string `json:"sender"`
	Content  string `json:"content"`
	DialogId string `json:"dialog_id"`
	TimeSend string `json:"time_send"`
}
