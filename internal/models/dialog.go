package models

type Dialog struct {
	DialogID        int    `json:"dialog_id"`
	UserIDOne       string `json:"user_id_1"`
	UserIDTwo       string `json:"user_id_2"`
	Avatar          string `json:"dialog_avatar"`
	LastMessage     string `json:"last_message"`
	UserOneUsername string `json:"user_one_username"`
	UserTwoUsername string `json:"user_two_username"`
}
