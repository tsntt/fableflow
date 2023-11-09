package memory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/account"
)

type AccountStore struct {
	db []account.Model
}

func NewAccountStore(as []account.Model) *AccountStore {
	return &AccountStore{
		db: as,
	}
}

func (store *AccountStore) New(a account.Model) (uuid.UUID, error) {
	a.ID = uuid.New()
	store.db = append(store.db, a)

	return a.ID, nil
}

func (store *AccountStore) GetByID(id uuid.UUID) (*account.Model, error) {
	var acc *account.Model
	for _, a := range store.db {
		if a.ID == id {
			acc = &a
		}
	}

	if acc == (&account.Model{}) {
		return nil, errors.New("transfer not fount")
	}

	return acc, nil
}

func (store *AccountStore) Update(a account.Model) error {
	idx := -1

	for i, acc := range store.db {
		if acc.ID == a.ID {
			idx = i
		}
	}

	if idx == -1 {
		return errors.New("transfer not fount")
	}

	store.db[idx] = a

	return nil
}

func (store *AccountStore) Delete(id uuid.UUID) error {
	idx := -1

	for i, acc := range store.db {
		if acc.ID == id {
			idx = i
		}
	}

	if idx == -1 {
		return errors.New("account not fount")
	}

	store.db = append(store.db[:idx], store.db[idx+1:]...)

	return nil
}
