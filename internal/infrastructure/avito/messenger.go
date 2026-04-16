package avito

import (
	"context"
	"prompter/internal/domain/entity"
	"prompter/internal/domain/port"
)

// Messenger реализует port.Messenger для Avito.
// Здесь живут все HTTP-вызовы к api.avito.ru/messenger/v3/*.
type Messenger struct{}

func NewMessenger() *Messenger {
	return &Messenger{}
}

func (m *Messenger) Provider() entity.Provider {
	return entity.ProviderAvito
}

func (m *Messenger) ListChats(ctx context.Context, token string, accountID string) ([]port.Chat, error) {
	panic("not implemented")
}

func (m *Messenger) GetMessages(ctx context.Context, token string, accountID string, chatID string) ([]port.ChatMessage, error) {
	panic("not implemented")
}

func (m *Messenger) SendMessage(ctx context.Context, token string, accountID string, chatID string, body string) (string, error) {
	panic("not implemented")
}

func (m *Messenger) GetVacancyIDForChat(ctx context.Context, token string, accountID string, chatID string) (int64, error) {
	panic("not implemented")
}
