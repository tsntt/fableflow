package bank

import (
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(i Model) (uuid.UUID, error)
	GetByDomain(domain Domain) (*Model, error)
	GetByHash(hash string) (*Model, error)
	Activate(i *Model) error
}

type Model struct {
	ID        uuid.UUID
	Domain    Domain
	Email     Email
	Hash      string
	CreatedAt time.Time
}

func New(d Domain, e Email) Model {
	return Model{
		ID:        uuid.New(),
		Domain:    d,
		Email:     e,
		CreatedAt: time.Now(),
	}
}

type Domain string

func NewDomain(domain string) (Domain, error) {
	// validade domain
	return Domain(domain), nil
}

func (d Domain) String() string { return string(d) }

type Email string

func NewEmail(e string) (Email, error) {
	// validade email
	return Email(e), nil
}

func (e Email) String() string { return string(e) }
