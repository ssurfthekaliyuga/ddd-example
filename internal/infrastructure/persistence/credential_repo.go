package persistence

import (
	"context"
	"fmt"
	"prompter/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

// CredentialModel — GORM-модель для таблицы credentials.
// Отдельная от Company таблица: credentials живут своим lifecycle (рефреш каждые N часов),
// а Company меняется редко. Это устраняет конфликты concurrent updates.
type CredentialModel struct {
	ID                    int64  `gorm:"primaryKey"`
	CompanyID             int64  `gorm:"index"`
	Provider              string `gorm:"index"`
	ClientID              string
	ClientSecretEncrypted string
	AccessToken           string
	ExpiresAt             time.Time
	RefreshedAt           time.Time
}

func (CredentialModel) TableName() string {
	return "credentials"
}

// CredentialRepo реализует port.CredentialStore через GORM + MySQL.
type CredentialRepo struct {
	db *gorm.DB
}

func NewCredentialRepo(db *gorm.DB) *CredentialRepo {
	return &CredentialRepo{db: db}
}

func (r *CredentialRepo) GetByCompanyAndProvider(ctx context.Context, companyID int64, provider entity.Provider) (*entity.Credential, error) {
	var model CredentialModel
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND provider = ?", companyID, string(provider)).
		First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("find credential: %w", err)
	}
	cred := toDomain(model)
	return &cred, nil
}

func (r *CredentialRepo) Save(ctx context.Context, cred *entity.Credential) error {
	model := toModel(*cred)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("save credential: %w", err)
	}
	cred.ID = model.ID
	return nil
}

func (r *CredentialRepo) ListExpiring(ctx context.Context) ([]entity.Credential, error) {
	var models []CredentialModel
	threshold := time.Now().Add(10 * time.Minute)
	err := r.db.WithContext(ctx).
		Where("expires_at < ? OR access_token = ''", threshold).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("list expiring: %w", err)
	}

	result := make([]entity.Credential, len(models))
	for i, m := range models {
		result[i] = toDomain(m)
	}
	return result, nil
}

func toDomain(m CredentialModel) entity.Credential {
	return entity.Credential{
		ID:                    m.ID,
		CompanyID:             m.CompanyID,
		Provider:              entity.Provider(m.Provider),
		ClientID:              m.ClientID,
		ClientSecretEncrypted: m.ClientSecretEncrypted,
		AccessToken:           m.AccessToken,
		ExpiresAt:             m.ExpiresAt,
		RefreshedAt:           m.RefreshedAt,
	}
}

func toModel(c entity.Credential) CredentialModel {
	return CredentialModel{
		ID:                    c.ID,
		CompanyID:             c.CompanyID,
		Provider:              string(c.Provider),
		ClientID:              c.ClientID,
		ClientSecretEncrypted: c.ClientSecretEncrypted,
		AccessToken:           c.AccessToken,
		ExpiresAt:             c.ExpiresAt,
		RefreshedAt:           c.RefreshedAt,
	}
}
