package accounts

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/account"
)

type AccountService struct {
	accountRepository account.Repository
}

func NewAccountService(ar account.Repository) *AccountService {
	return &AccountService{
		accountRepository: ar,
	}
}

type NewAccountReqDTO struct {
	BankID         uuid.UUID `json:"bank_id"`
	InitialDeposit float64   `json:"initial_deposit"`
}

type NewAccountRespDTO struct {
	ID uuid.UUID `json:"id"`
}

type AccountRespDTO struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}

func (s *AccountService) Create(ar NewAccountReqDTO) (*NewAccountRespDTO, error) {
	if ar.InitialDeposit < 0 {
		return nil, errors.New("account cannot be created with negative balance")
	}

	newAcc := account.New(ar.BankID, 0.0)

	id, err := s.accountRepository.Create(newAcc)
	if err != nil {
		return nil, err
	}

	return &NewAccountRespDTO{
		ID: id,
	}, nil
}

func (s *AccountService) GetByID(id uuid.UUID) (*account.Model, error) {
	acc, err := s.accountRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) Update(a account.Model) error {
	err := s.accountRepository.UpdateBalance(a.ID, a.Balance)

	return err
}
