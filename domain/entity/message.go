package entity

import "time"

// IncomingMessage — входящее сообщение от кандидата (бывший Call).
type IncomingMessage struct {
	ID          int64
	CreatedAt   time.Time
	InterviewID int64
	ExternalID  string // AvitoMessageId / HH message id
	AuthorID    int64
	MessageType string // text, image, etc.
	Body        string
	IsProcessed bool
	HideFromBroker bool
}

// OutgoingMessage — исходящее сообщение для отправки (бывший Cue).
type OutgoingMessage struct {
	ID          int64
	CreatedAt   time.Time
	InterviewID int64
	Body        string
	ExternalID  string // заполняется после отправки
	IsExternal  bool   // создано человеком, а не ботом
	DontSend    bool
	HideFromBroker bool
}
