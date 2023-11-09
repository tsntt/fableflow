package transfers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/transfer"
)

type TransferService struct {
	transferRepository transfer.Repository
}

func NewTransferService(TransferRepo transfer.Repository) *TransferService {
	return &TransferService{
		transferRepository: TransferRepo,
	}
}

type NewTransferReqDTO struct {
	Receiver  uuid.UUID `json:"receiver"`
	Sender    uuid.UUID `json:"sender"`
	Amount    float64   `json:"amount"`
	Scheduled time.Time `json:"scheduled"`
}

type NewTransferRespDTO struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (s *TransferService) CreateTransfer(tr NewTransferReqDTO) (*transfer.Model, error) {
	val, err := transfer.NewValue(tr.Amount)
	if err != nil {
		return nil, err
	}

	scd, err := transfer.NewShedule(tr.Scheduled)
	if err != nil {
		return nil, err
	}

	t := transfer.New(tr.Receiver, tr.Sender, val, scd)

	id, err := s.transferRepository.Create(t)
	if err != nil {
		return nil, err
	}

	t.ID = id

	return &t, nil
}

func (s *TransferService) Approve(id uuid.UUID) (*transfer.Model, error) {
	tr, err := s.transferRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if tr.Status != transfer.Created && tr.Status != transfer.Pending {
		return nil, errors.New("transaction aready processed")
	}

	err = s.transferRepository.UpdateStatus(id, transfer.Approved)
	if err != nil {
		return nil, err
	}

	tr.Status = transfer.Approved

	return tr, nil
}

func (s *TransferService) Rejecte(id uuid.UUID) (*transfer.Model, error) {
	tr, err := s.transferRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if tr.Status != transfer.Created && tr.Status != transfer.Pending {
		return nil, errors.New("transaction aready processed")
	}

	err = s.transferRepository.UpdateStatus(id, transfer.Rejected)
	if err != nil {
		return nil, err
	}

	tr.Status = transfer.Rejected

	return tr, nil
}

func (s *TransferService) Schedule(id uuid.UUID) (*transfer.Model, error) {
	tr, err := s.transferRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if tr.Status != transfer.Created {
		return nil, errors.New("transaction aready processed")
	}

	err = s.transferRepository.UpdateStatus(id, transfer.Pending)
	if err != nil {
		return nil, err
	}

	tr.Status = transfer.Pending

	return tr, nil
}

func (s *TransferService) GetByID(id uuid.UUID) (*transfer.Model, error) {
	t, err := s.transferRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TransferService) GetByAccountAndPeriod(accID uuid.UUID, start, end time.Time) ([]transfer.Model, error) {
	dates, err := transfer.NewSearchDate(start, end)
	if err != nil {
		return nil, err
	}

	ts, err := s.transferRepository.GetByAccountAndPeriod(accID, dates)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TransferService) GetTodaysPending() ([]transfer.Model, error) {
	period, err := transfer.NewSearchDate(time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, err
	}

	ts, err := s.transferRepository.GetByPeriodAndStatus(period, transfer.Pending)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TransferService) Cancel(id uuid.UUID) error {
	tr, err := s.transferRepository.GetByID(id)
	if err != nil {
		return err
	}

	if tr.Status != transfer.Created && tr.Status != transfer.Pending {
		return errors.New("transaction aready processed, cannot be undone")
	}

	err = s.transferRepository.UpdateStatus(id, transfer.Canceled)

	return err
}
