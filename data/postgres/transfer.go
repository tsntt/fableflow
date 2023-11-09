package postgres

import (
	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/transfer"
)

type TransferStorage PostgresStorage

func NewTransferStorage(s *PostgresStorage) *TransferStorage {
	transferStore := new(TransferStorage)
	transferStore.db = s.db
	return transferStore
}

func (s *TransferStorage) Create(t transfer.Model) (uuid.UUID, error) {
	q := `INSERT INTO transfers(receiver, sender, amount, status, scheduled, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// allow sender or receiver to be null
	var sender any
	if (t.Sender == uuid.UUID{}) && (t.Receiver != uuid.UUID{}) {
		sender = nil
	} else {
		sender = t.Sender
	}

	var receiver any
	if (t.Receiver == uuid.UUID{}) && (t.Sender != uuid.UUID{}) {
		receiver = nil
	} else {
		receiver = t.Receiver
	}

	var id uuid.UUID

	err := s.db.QueryRow(q, receiver, sender, t.Amount, t.Status, t.Scheduled.Time(), &t.CreatedAt).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (s *TransferStorage) GetByID(id uuid.UUID) (*transfer.Model, error) {
	q := `SELECT * FROM transfers WHERE id=$1`

	tr := transfer.Model{}

	err := s.db.QueryRow(q, id).Scan(
		&tr.ID,
		&tr.Receiver,
		&tr.Sender,
		&tr.Amount,
		&tr.Status,
		&tr.Scheduled,
		&tr.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &tr, nil
}

// remove
func (s *TransferStorage) GetByAccount(accountID uuid.UUID) ([]transfer.Model, error) {
	q := `SELECT * FROM transfers WHERE receiver=$1 OR sender=$1`

	trs := []transfer.Model{}

	rows, err := s.db.Query(q, accountID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tr := transfer.Model{}

		if err := rows.Scan(
			&tr.ID,
			&tr.Receiver,
			&tr.Sender,
			&tr.Amount,
			&tr.Status,
			&tr.Scheduled,
			&tr.CreatedAt,
		); err != nil {
			return nil, err
		}

		trs = append(trs, tr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trs, nil
}

func (s *TransferStorage) GetByAccountAndPeriod(accountID uuid.UUID, d transfer.Period) ([]transfer.Model, error) {
	q := `SELECT * FROM transfers WHERE sender=$1 OR receiver=$1 AND scheduled > $2 AND scheduled < $3`

	trs := []transfer.Model{}

	rows, err := s.db.Query(q, accountID, d.Start, d.End)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tr := transfer.Model{}

		if err := rows.Scan(
			&tr.ID,
			&tr.Receiver,
			&tr.Sender,
			&tr.Amount,
			&tr.Status,
			&tr.Scheduled,
			&tr.CreatedAt,
		); err != nil {
			return nil, err
		}

		trs = append(trs, tr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trs, nil
}

func (s *TransferStorage) GetByPeriodAndStatus(p transfer.Period, status transfer.Status) ([]transfer.Model, error) {
	q := `SELECT * FROM transfers WHERE scheduled > $1 AND scheduled < $2 AND status=$3`

	trs := []transfer.Model{}

	rows, err := s.db.Query(q, p.Start, p.End, status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tr := transfer.Model{}

		if err := rows.Scan(
			&tr.ID,
			&tr.Receiver,
			&tr.Sender,
			&tr.Amount,
			&tr.Status,
			&tr.Scheduled,
			&tr.CreatedAt,
		); err != nil {
			return nil, err
		}

		trs = append(trs, tr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trs, nil
}

func (s *TransferStorage) UpdateStatus(id uuid.UUID, status transfer.Status) error {
	q := `UPDATE transfers SET status=$2 WHERE id=$1 RETURNING *`
	_, err := s.db.Exec(q, id, status)

	return err
}
