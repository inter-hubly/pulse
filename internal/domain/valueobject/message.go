package valueobject

type Message struct {
	owner   bool   `json:"owner"`
	message string `json:"message"`
	time    int64  `json:"time"`
}

func NewMessage(owner bool, message string, time int64) *Message {
	return &Message{
		owner:   owner,
		message: message,
		time:    time,
	}
}

func (m *Message) GetMessage() string {
	return m.message
}
