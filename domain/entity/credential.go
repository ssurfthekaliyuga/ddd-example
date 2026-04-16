package entity

import "time"

// Provider — внешняя площадка, с которой мы интегрируемся.
type Provider string

const (
	ProviderAvito Provider = "avito"
	ProviderHH    Provider = "hh"
)

// Credential — учётные данные для внешнего API.
// Это НЕ доменная сущность в классическом смысле, а value object инфраструктурного слоя.
// Но мы описываем его в domain, потому что порты (интерфейсы) оперируют им.
type Credential struct {
	ID        int64
	CompanyID int64
	Provider  Provider

	// Зашифрованные данные для получения токена (client_id, encrypted secret и т.д.)
	ClientID              string
	ClientSecretEncrypted string

	// Текущий живой токен
	AccessToken string
	ExpiresAt   time.Time
	RefreshedAt time.Time
}

// IsExpired проверяет, нужно ли обновить токен.
// Бизнес-правило: обновляем заранее, за 10 минут до истечения.
func (c Credential) IsExpired() bool {
	return time.Now().After(c.ExpiresAt.Add(-10 * time.Minute))
}
