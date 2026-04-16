package usecase

import (
	"context"
	"fmt"
	"hr-bot-ddd-example/internal/domain/entity"
	"hr-bot-ddd-example/internal/domain/port"
	"time"

	"go.uber.org/zap"
)

// RefreshTokens — use case: обновить истекающие токены всех площадок.
// Вызывается кроном (infrastructure/scheduler), сам крон ничего не знает про токены.
type RefreshTokens struct {
	credStore  port.CredentialStore
	refreshers port.TokenRefreshersRegistry
	decryptKey string
	logger     *zap.SugaredLogger
}

func NewRefreshTokens(
	credStore port.CredentialStore,
	refreshers port.TokenRefreshersRegistry,
	decryptKey string,
	logger *zap.SugaredLogger,
) *RefreshTokens {
	return &RefreshTokens{
		credStore:  credStore,
		refreshers: refreshers,
		decryptKey: decryptKey,
		logger:     logger,
	}
}

// Execute находит все credentials с истекающим токеном и обновляет их.
func (uc *RefreshTokens) Execute(ctx context.Context) error {
	expiring, err := uc.credStore.ListExpiring(ctx)
	if err != nil {
		return fmt.Errorf("list expiring credentials: %w", err)
	}

	for _, cred := range expiring {
		if err := uc.refreshOne(ctx, cred); err != nil {
			// Ошибка одного не останавливает остальные
			uc.logger.Errorw("token refresh failed",
				"company_id", cred.CompanyID,
				"provider", cred.Provider,
				"error", err,
			)
			continue
		}
	}

	return nil
}

func (uc *RefreshTokens) refreshOne(ctx context.Context, cred entity.Credential) error {
	refresher, err := uc.refreshers.Get(cred.Provider)
	if err != nil {
		return fmt.Errorf("no refresher for provider %s", cred.Provider)
	}

	token, expiresIn, err := refresher.Refresh(ctx, cred, uc.decryptKey)
	if err != nil {
		return fmt.Errorf("refresh %s token for company %d: %w", cred.Provider, cred.CompanyID, err)
	}

	cred.AccessToken = token
	cred.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	cred.RefreshedAt = time.Now()

	if err := uc.credStore.Save(ctx, &cred); err != nil {
		return fmt.Errorf("save refreshed credential: %w", err)
	}

	uc.logger.Infow("token refreshed",
		"company_id", cred.CompanyID,
		"provider", cred.Provider,
		"expires_at", cred.ExpiresAt,
	)

	return nil
}
