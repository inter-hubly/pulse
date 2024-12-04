package valueobject

type Message struct {
	ToId    string `json:"toId"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

func NewMessage(toId, message string, time int64) *Message {
	return &Message{
		ToId:    toId,
		Message: message,
		Time:    time,
	}
}
