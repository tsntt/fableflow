package postgres

import (
	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/account"
)

type AccountStorage PostgresStorage

func NewAccountStorage(s *PostgresStorage) *AccountStorage {
	accountStore := new(AccountStorage)
	accountStore.db = s.db
	return accountStore
}

func (s *AccountStorage) Create(a account.Model) (uuid.UUID, error) {
	q := `INSERT INTO accounts(bank_id, balance, created_at)
		VALUES ($1, $2, $3) RETURNING id`

	var id uuid.UUID

	err := s.db.QueryRow(q, a.BankID, a.Balance, a.CreatedAt).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (s *AccountStorage) GetByID(id uuid.UUID) (*account.Model, error) {
	q := `SELECT * FROM accounts WHERE id=$1`

	acc := account.Model{}

	err := s.db.QueryRow(q, id).Scan(
		&acc.ID,
		&acc.BankID,
		&acc.Balance,
		&acc.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &acc, nil
}
func (s *AccountStorage) UpdateBalance(id uuid.UUID, balance float64) error {
	q := `UPDATE accounts SET balance=$2 WHERE id=$1 RETURNING *`
	_, err := s.db.Exec(q, id, balance)

	return err
}
