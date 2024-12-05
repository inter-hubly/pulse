package valueobject

type Message struct {
	ToId    string `json:"toId"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
	IsOwner bool   `json:"isOwner"`
}

func NewMessage(toId, message string, time int64, isOwner bool) *Message {
	return &Message{
		ToId:    toId,
		Message: message,
		Time:    time,
		IsOwner: isOwner,
	}
}
