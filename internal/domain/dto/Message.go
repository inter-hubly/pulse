package dto

type Message struct {
	Username string `json:"username"`
	ToId     string `json:"toId"`
	Message  string `json:"message"`
	IsOwner  bool   `json:"isOwner"`
}
