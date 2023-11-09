package memory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/transfer"
)

type TransferStore struct {
	db []transfer.Model
}

func NewTransferStore(ts []transfer.Model) *TransferStore {
	return &TransferStore{
		db: ts,
	}
}

func (store *TransferStore) New(t transfer.Model) (uuid.UUID, error) {
	t.ID = uuid.New()
	store.db = append(store.db, t)

	return t.ID, nil
}

func (store *TransferStore) GetByID(id uuid.UUID) (*transfer.Model, error) {
	var tr *transfer.Model
	for _, t := range store.db {
		if t.ID == id {
			tr = &t
		}
	}

	if tr == (&transfer.Model{}) {
		return nil, errors.New("transfer not fount")
	}

	return tr, nil
}

func (store *TransferStore) GetByAccount(accountID uuid.UUID) ([]transfer.Model, error) {
	var transfers []transfer.Model

	for _, transfer := range store.db {
		if transfer.Receiver == accountID || transfer.Sender == accountID {
			transfers = append(transfers, transfer)
		}
	}

	return transfers, nil
}

func (store *TransferStore) GetByDate(accountID uuid.UUID, d transfer.Period) ([]transfer.Model, error) {
	var transfers []transfer.Model

	for _, transfer := range store.db {
		if (transfer.Receiver == accountID || transfer.Sender == accountID) &&
			transfer.CreatedAt.Before(d.End) && transfer.CreatedAt.After(d.Start) {
			transfers = append(transfers, transfer)
		}
	}

	return transfers, nil
}

func (store *TransferStore) GetAllByDate(d transfer.Period, status transfer.Status) ([]transfer.Model, error) {
	var transfers []transfer.Model

	for _, transfer := range store.db {
		if transfer.Status == status && transfer.CreatedAt.Before(d.End) && transfer.CreatedAt.After(d.Start) {
			transfers = append(transfers, transfer)
		}
	}

	return transfers, nil
}

func (store *TransferStore) Update(t *transfer.Model) (*transfer.Model, error) {
	idx := -1

	for i, tr := range store.db {
		if tr.ID == t.ID {
			idx = i
		}
	}

	if idx == -1 {
		return nil, errors.New("transfer not fount")
	}

	store.db[idx].Status = t.Status

	return &store.db[idx], nil
}
