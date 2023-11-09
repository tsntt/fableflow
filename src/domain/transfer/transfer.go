package transfer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNegativeValue = errors.New("cannot transfer a negative value")
	ErrZeroValue     = errors.New("can't transfer nothing")
)

type Repository interface {
	Create(t Model) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*Model, error)
	GetByPeriodAndStatus(p Period, status Status) ([]Model, error)
	GetByAccountAndPeriod(accountID uuid.UUID, p Period) ([]Model, error)
	UpdateStatus(id uuid.UUID, status Status) error
}

// TODO: make DTO for this type
type Model struct {
	ID        uuid.UUID `json:"id"`
	Receiver  uuid.UUID `json:"receiver"`
	Sender    uuid.UUID `json:"sender"`
	Amount    Amount    `json:"amount"`
	Status    Status    `json:"status"`
	Scheduled Scheduled `json:"scheduled"`
	CreatedAt time.Time `json:"created_at"`
}

func New(sender, receiver uuid.UUID, val Amount, date Scheduled) Model {
	return Model{
		Receiver:  sender,
		Sender:    receiver,
		Amount:    val,
		Status:    Created,
		Scheduled: date,
		CreatedAt: time.Now().UTC(),
	}
}

type Scheduled time.Time

func NewShedule(t time.Time) (Scheduled, error) {
	now := time.Now()
	aYear := now.Add(time.Hour * 24 * 365)

	if t.After(aYear) {
		return Scheduled{}, errors.New("cannot shedule more than a year in the future")
	}

	if t.Before(now) {
		return Scheduled(now), nil
	}

	return Scheduled(t), nil
}

func (s Scheduled) Time() time.Time {
	return time.Time(s)
}

type Status string

const (
	Created  Status = "created"
	Canceled Status = "canceled"
	Pending  Status = "pending"
	Approved Status = "approved"
	Rejected Status = "rejected"
)

func (m *Model) Update(b bool) {
	if b {
		m.Status = Approved
	} else {
		m.Status = Rejected
	}
}

type Amount float64

func NewValue(val float64) (Amount, error) {
	if val < 0 {
		return 0, ErrNegativeValue
	}

	if val == 0 {
		return 0, ErrZeroValue
	}

	return Amount(val), nil
}
