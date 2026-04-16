package port

import (
	"context"
	"prompter/ddd_reference/domain/entity"
)

// TokenRefresher — порт для получения нового токена от внешней площадки.
// Каждая площадка (Avito, HH) реализует свой — это infrastructure.
type TokenRefresher interface {
	// Provider возвращает площадку, которую обслуживает этот refresher.
	Provider() entity.Provider

	// Refresh получает новый access token, используя credentials.
	// Возвращает новый токен и время его истечения.
	Refresh(ctx context.Context, cred entity.Credential, decryptionKey string) (token string, expiresIn int, err error)
}

type TokenRefreshersRegistry interface {
	Get(entity.Provider) (TokenRefresher, error)
}
