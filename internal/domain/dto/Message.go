package dto

type Message struct {
	ProfileName string `json:"profileName,omitempty"`
	ToPhone     string `json:"toPhone"`
	Message     string `json:"message"`
	IsOwner     bool   `json:"isOwner"`
}

type Template struct {
	ToPhone  string `json:"toPhone"`
	Name     string `json:"name"`
	Language string `json:"language"`
}
