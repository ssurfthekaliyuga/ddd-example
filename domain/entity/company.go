package entity

import "time"

// Company — доменная сущность: аккаунт работодателя на площадке.
// Никаких токенов, секретов, инфраструктурных полей.
type Company struct {
	ID        int64
	CreatedAt time.Time
	Title     string

	// Идентификаторы на внешних площадках — это бизнес-факт ("мы зарегистрированы там-то"),
	// а не инфраструктура. Домен знает ЧТО у нас есть аккаунт, но не знает КАК мы авторизуемся.
	AvitoAccountId string
	HHEmployerId   string

	CheckChats bool // опрашивать ли чаты этой компании
}

func (c *Company) GetAccountID(provider Provider) string {
	switch provider {
	case ProviderAvito:
		return c.AvitoAccountId
	case ProviderHH:
		return c.HHEmployerId
	default:
		return ""
	}
}
