package scheduler

import (
	"context"
	"hr-bot-ddd-example/internal/application/usecase"
	"hr-bot-ddd-example/internal/domain/entity"
	"time"

	"go.uber.org/zap"
)

// StartTickers запускает все фоновые циклы.
// Scheduler — чистая infrastructure: он знает КОГДА запускать, но не знает ЧТО делать.
// Вся логика — в use cases.
func StartTickers(
	refreshTokens *usecase.RefreshTokens,
	checkChats *usecase.CheckChats,
	logger *zap.SugaredLogger,
) {
	// Refresh токенов — раз в час (с запасом до истечения)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		// Сразу при старте
		if err := refreshTokens.Execute(context.Background()); err != nil {
			logger.Errorw("initial token refresh failed", "error", err)
		}

		for range ticker.C {
			if err := refreshTokens.Execute(context.Background()); err != nil {
				logger.Errorw("token refresh tick failed", "error", err)
			}
		}
	}()

	// Check chats — раз в минуту для каждой площадки
	providers := []entity.Provider{entity.ProviderAvito, entity.ProviderHH}

	for _, p := range providers {
		provider := p
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()

			for range ticker.C {
				if err := checkChats.Execute(context.Background(), provider); err != nil {
					logger.Errorw("check chats tick failed",
						"provider", provider,
						"error", err,
					)
				}
			}
		}()
	}
}
