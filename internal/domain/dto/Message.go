package dto

type Message struct {
	Username string `json:"username"`
	ToNumber string `json:"toNumber"`
	Message  string `json:"message"`
	IsOwner  bool   `json:"isOwner"`
}

type Template struct {
	OwnerId  string `json:"ownerId"`
	ToNumber string `json:"toNumber"`
	Name     string `json:"name"`
	Language string `json:"language"`
}
