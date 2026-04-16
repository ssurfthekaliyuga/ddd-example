package port

import (
	"context"
	"hr-bot-ddd-example/internal/domain/entity"
)

// CredentialStore — порт для хранения/получения учётных данных.
// Реализация — infrastructure/persistence.
type CredentialStore interface {
	// GetByCompanyAndProvider возвращает credential для конкретной компании и площадки.
	GetByCompanyAndProvider(ctx context.Context, companyID int64, provider entity.Provider) (*entity.Credential, error)

	// Save сохраняет credential (в т.ч. обновлённый токен).
	Save(ctx context.Context, cred *entity.Credential) error

	// ListExpiring возвращает все credentials, которым скоро нужен refresh.
	ListExpiring(ctx context.Context) ([]entity.Credential, error)
}
