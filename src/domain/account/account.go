package account

import (
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(a Model) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*Model, error)
	UpdateBalance(id uuid.UUID, balance float64) error
}

type Model struct {
	ID        uuid.UUID `json:"id"`
	BankID    uuid.UUID `json:"-"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"date"`
}

func New(bankID uuid.UUID, value float64) Model {
	return Model{
		BankID:    bankID,
		Balance:   value,
		CreatedAt: time.Now(),
	}
}
