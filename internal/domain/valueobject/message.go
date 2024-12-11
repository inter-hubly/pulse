package valueobject

import "time"

// Message struct representando uma mensagem
type Message struct {
	MessageId   string `json:"messageId"`
	ToPhone     string `json:"toPhone,omitempty"`
	ProfileName string `json:"profileName,omitempty"`
	Message     string `json:"message"`
	Time        int64  `json:"time"`
	IsOwner     bool   `json:"isOwner"`
}

type Option func(*Message)

func NewMessage(options ...Option) *Message {
	msg := &Message{
		Time: time.Now().Unix(), // Valor padrão
	}
	for _, opt := range options {
		opt(msg)
	}
	return msg
}

func WithMessageId(messageId string) Option {
	return func(msg *Message) {
		msg.MessageId = messageId
	}
}

func WithToPhone(toPhone string) Option {
	return func(m *Message) {
		m.ToPhone = toPhone
	}
}

func WithProfileName(profileName string) Option {
	return func(m *Message) {
		m.ProfileName = profileName
	}
}

func WithMessage(message string) Option {
	return func(m *Message) {
		m.Message = message
	}
}

func WithIsOwner(isOwner bool) Option {
	return func(m *Message) {
		m.IsOwner = isOwner
	}
}
