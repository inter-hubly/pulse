package valueobject

import "time"

// Message struct representando uma mensagem
type Message struct {
	ToId        string `json:"toId,omitempty"`
	ProfileName string `json:"profileName,omitempty"`
	Message     string `json:"message"`
	Time        int64  `json:"time"`
	IsOwner     bool   `json:"isOwner"`
}

// Option é o tipo de função que configura o objeto Message
type Option func(*Message)

// NewMessage cria um novo objeto Message com as opções fornecidas
func NewMessage(options ...Option) *Message {
	msg := &Message{
		Time: time.Now().Unix(), // Valor padrão
	}
	for _, opt := range options {
		opt(msg)
	}
	return msg
}

// WithToId configura o campo ToId
func WithToId(toId string) Option {
	return func(m *Message) {
		m.ToId = toId
	}
}

// WithProfileName configura o campo ProfileName
func WithProfileName(profileName string) Option {
	return func(m *Message) {
		m.ProfileName = profileName
	}
}

// WithMessage configura o campo Message
func WithMessage(message string) Option {
	return func(m *Message) {
		m.Message = message
	}
}

// WithIsOwner configura o campo IsOwner
func WithIsOwner(isOwner bool) Option {
	return func(m *Message) {
		m.IsOwner = isOwner
	}
}
