package banks

import (
	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/bank"
	"github.com/tsntt/fableflow/src/util"
)

type BankService struct {
	bankRepository bank.Repository
}

func NewBankService(br bank.Repository) *BankService {
	return &BankService{
		bankRepository: br,
	}
}

type NewBankReqDTO struct {
	Domain string `json:"domain"`
	Email  string `json:"email"`
}

type NewBankRespDTO struct {
	ID   *uuid.UUID `json:"id"`
	Hash string     `json:"-"`
}

func (s *BankService) Create(r NewBankReqDTO) (*NewBankRespDTO, error) {
	domain, err := bank.NewDomain(r.Domain)
	if err != nil {
		return nil, err
	}

	email, err := bank.NewEmail(r.Email)
	if err != nil {
		return nil, err
	}

	bk := bank.New(domain, email)
	bk.Hash = util.RandHash()

	id, err := s.bankRepository.Create(bk)
	if err != nil {
		return nil, err
	}

	return &NewBankRespDTO{
		ID:   &id,
		Hash: bk.Hash,
	}, err
}

func (s *BankService) Get(d string) (*bank.Model, error) {
	domain, err := bank.NewDomain(d)
	if err != nil {
		return nil, err
	}

	bk, err := s.bankRepository.GetByDomain(domain)
	if err != nil {
		return nil, err
	}

	return bk, nil
}

func (s *BankService) Activate(hash string) error {
	bk, err := s.bankRepository.GetByHash(hash)
	if err != nil {
		return err
	}

	err = s.bankRepository.Activate(bk)

	return err
}
