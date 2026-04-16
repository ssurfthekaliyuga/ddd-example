package port

import (
	"context"
	"prompter/ddd_reference/domain/entity"
)

// Chat — результат опроса чатов на площадке.
type Chat struct {
	ExternalChatID string
	UpdatedAtTs    int64
	Messages       []ChatMessage
}

type ChatMessage struct {
	ExternalID  string
	AuthorID    int64
	Direction   string // "in" / "out"
	Body        string
	MessageType string
	CreatedAtTs int64
}

// Messenger — порт для взаимодействия с мессенджером площадки (Avito, HH).
// Реализация — infrastructure/avito, infrastructure/hh.
type Messenger interface {
	Provider() entity.Provider

	// ListChats возвращает активные чаты для компании.
	ListChats(ctx context.Context, token string, accountID string) ([]Chat, error)

	// GetMessages возвращает сообщения конкретного чата.
	GetMessages(ctx context.Context, token string, accountID string, chatID string) ([]ChatMessage, error)

	// SendMessage отправляет сообщение в чат, возвращает внешний ID сообщения.
	SendMessage(ctx context.Context, token string, accountID string, chatID string, body string) (externalMessageID string, err error)

	// GetVacancyIDForChat определяет, к какой вакансии относится чат.
	GetVacancyIDForChat(ctx context.Context, token string, accountID string, chatID string) (externalVacancyID int64, err error)
}

// MessengerRegistry — реестр мессенджеров по площадкам.
// Позволяет use case'ам не зависеть от конкретной площадки.
type MessengerRegistry interface {
	Get(provider entity.Provider) (Messenger, error)
}
