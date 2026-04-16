package entity

import "time"

type Interview struct {
	ID            int64
	CreatedAt     time.Time
	IsBotDisabled bool
	CandidateID   int64
	VacancyID     int64
	BrokerID      *int64 // переопределение брокеров, nil = использовать из Vacancy
	Class         string
	ChatID        string // абстрактный ID чата (Avito, HH, etc.)
	ChatUpdatedAt time.Time
}

type Candidate struct {
	ID        int64
	CreatedAt time.Time
	IsStaff   bool
	Name      string

	// Внешние ID — по площадкам.
	// В отличие от Company, кандидат привязан к конкретной площадке.
	AvitoID  int64
	AvitoURL string
	HHID     string
}

type Vacancy struct {
	ID                int64
	CreatedAt         time.Time
	Title             string
	CompanyID         int64
	NewInterviewClass string
	BrokerIdsJson     string

	// Внешние ID вакансий на площадках
	AvitoVacancyID int64
	HHVacancyID    string
}
