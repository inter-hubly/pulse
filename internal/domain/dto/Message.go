package dto

type Message struct {
	Username string `json:"username"`
	ToId     string `json:"toId"`
	Message  string `json:"message"`
	IsOwner  bool   `json:"isOwner"`
}

type Template struct {
	OwnerId  string `json:"ownerId"`
	ToId     string `json:"toId"`
	Name     string `json:"name"`
	Language string `json:"language"`
}
