package valueobject

type Message struct {
	owner   string `json:"owner"`
	message string `json:"message"`
	time    int64  `json:"time"`
}

func NewMessage(owner, message string, time int64) *Message {
	return &Message{
		owner:   owner,
		message: message,
		time:    time,
	}
}

func (m *Message) GetMessage() string {
	return m.message
}

func (m *Message) GetOwner() string {
	return m.owner
}
