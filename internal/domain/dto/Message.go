package dto

type Message struct {
	Username string `json:"username"`
	ToPhone  string `json:"toPhone"`
	Message  string `json:"message"`
	IsOwner  bool   `json:"isOwner"`
}

type Template struct {
	OwnerId  string `json:"ownerId"`
	ToPhone  string `json:"toPhone"`
	Name     string `json:"name"`
	Language string `json:"language"`
}
