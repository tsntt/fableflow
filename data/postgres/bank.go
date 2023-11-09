package postgres

import (
	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/bank"
)

type BankStorage PostgresStorage

func NewBankStorage(s *PostgresStorage) *BankStorage {
	bankStorage := new(BankStorage)
	bankStorage.db = s.db
	return bankStorage
}

func (s *BankStorage) Create(i bank.Model) (uuid.UUID, error) {
	q := `INSERT INTO tmp_banks(domain, email, hash, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id uuid.UUID
	err := s.db.QueryRow(q, &i.Domain, &i.Email, &i.Hash, &i.CreatedAt).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (s *BankStorage) GetByDomain(domain bank.Domain) (*bank.Model, error) {
	q := `SELECT * FROM banks WHERE domain=$1`

	bk := bank.Model{}

	err := s.db.QueryRow(q, domain.String()).Scan(
		&bk.ID,
		&bk.Domain,
		&bk.Email,
		&bk.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &bk, nil
}

func (s *BankStorage) GetByHash(hash string) (*bank.Model, error) {
	q := `SELECT * FROM tmp_banks WHERE hash=$1`

	bk := bank.Model{}

	err := s.db.QueryRow(q, hash).Scan(
		&bk.ID,
		&bk.Domain,
		&bk.Email,
		&bk.Hash,
		&bk.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &bk, nil
}

func (s *BankStorage) Activate(b *bank.Model) error {
	q := `INSERT INTO banks(id, domain, email, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(q, b.ID, b.Domain, b.Email, b.CreatedAt)

	return err
}
