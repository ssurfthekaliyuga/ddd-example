package usecase

import (
	"context"
	"fmt"
	"prompter/ddd_reference/domain/entity"
	"prompter/ddd_reference/domain/port"

	"go.uber.org/zap"
)

// CheckChats — use case: опросить чаты площадки, синхронизировать сообщения в БД.
// Аналог текущего checkChatsTick, но без привязки к конкретной площадке.
type CheckChats struct {
	credStore   port.CredentialStore
	messengers  port.MessengerRegistry
	companyRepo CompanyRepo
	chatRepo    ChatRepo
	logger      *zap.SugaredLogger
}

// CompanyRepo — минимальный интерфейс, нужный этому use case.
// Определяем здесь, а не в port/, потому что это ad-hoc зависимость use case.
type CompanyRepo interface {
	ListActive(ctx context.Context) ([]entity.Company, error)
}

// ChatRepo — репозиторий для Interview/Call/Cue/Candidate.
type ChatRepo interface {
	FindOrCreateInterview(ctx context.Context, companyID int64, chatID string, vacancyID int64, candidateName string) (*entity.Interview, error)
	SaveIncomingMessage(ctx context.Context, msg *entity.IncomingMessage) error
	SaveOutgoingMessage(ctx context.Context, msg *entity.OutgoingMessage) error
	FindVacancyByExternalID(ctx context.Context, provider entity.Provider, externalID int64) (*entity.Vacancy, error)
}

func NewCheckChats(
	credStore port.CredentialStore,
	messengers port.MessengerRegistry,
	companyRepo CompanyRepo,
	chatRepo ChatRepo,
	logger *zap.SugaredLogger,
) *CheckChats {
	return &CheckChats{
		credStore:   credStore,
		messengers:  messengers,
		companyRepo: companyRepo,
		chatRepo:    chatRepo,
		logger:      logger,
	}
}

// Execute опрашивает все площадки для всех активных компаний.
func (uc *CheckChats) Execute(ctx context.Context, provider entity.Provider) error {
	companies, err := uc.companyRepo.ListActive(ctx)
	if err != nil {
		return fmt.Errorf("list active companies: %w", err)
	}

	messenger, err := uc.messengers.Get(provider)
	if err != nil {
		return fmt.Errorf("get messenger for %s: %w", provider, err)
	}

	for _, company := range companies {
		if err := uc.checkCompanyChats(ctx, company, messenger); err != nil {
			uc.logger.Errorw("check chats failed",
				"company_id", company.ID,
				"provider", provider,
				"error", err,
			)
		}
	}

	return nil
}

func (uc *CheckChats) checkCompanyChats(ctx context.Context, company entity.Company, messenger port.Messenger) error {
	cred, err := uc.credStore.GetByCompanyAndProvider(ctx, company.ID, messenger.Provider())
	if err != nil {
		return fmt.Errorf("get credential: %w", err)
	}

	accountID := company.GetAccountID(messenger.Provider())

	chats, err := messenger.ListChats(ctx, cred.AccessToken, accountID)
	if err != nil {
		return fmt.Errorf("list chats: %w", err)
	}

	for _, chat := range chats {
		vacancyExtID, err := messenger.GetVacancyIDForChat(ctx, cred.AccessToken, accountID, chat.ExternalChatID)
		if err != nil {
			uc.logger.Warnw("skip chat: cannot get vacancy", "chat_id", chat.ExternalChatID, "error", err)
			continue
		}

		vacancy, err := uc.chatRepo.FindVacancyByExternalID(ctx, messenger.Provider(), vacancyExtID)
		if err != nil {
			continue
		}

		interview, err := uc.chatRepo.FindOrCreateInterview(ctx, company.ID, chat.ExternalChatID, vacancy.ID, "")
		if err != nil {
			uc.logger.Errorw("find/create interview", "error", err)
			continue
		}

		for _, msg := range chat.Messages {
			if msg.Direction == "in" {
				_ = uc.chatRepo.SaveIncomingMessage(ctx, &entity.IncomingMessage{
					InterviewID: interview.ID,
					ExternalID:  msg.ExternalID,
					AuthorID:    msg.AuthorID,
					MessageType: msg.MessageType,
					Body:        msg.Body,
				})
			} else {
				_ = uc.chatRepo.SaveOutgoingMessage(ctx, &entity.OutgoingMessage{
					InterviewID: interview.ID,
					ExternalID:  msg.ExternalID,
					Body:        msg.Body,
					IsExternal:  true, // сообщение от человека, обнаруженное при синхронизации
				})
			}
		}
	}

	return nil
}
