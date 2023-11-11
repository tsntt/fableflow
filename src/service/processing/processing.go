package processing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/src/domain/account"
	"github.com/tsntt/fableflow/src/domain/transfer"
	"github.com/tsntt/fableflow/src/service/accounts"
	"github.com/tsntt/fableflow/src/service/transfers"
)

type ProcessingService struct {
	accountService  *accounts.AccountService
	transferService *transfers.TransferService
}

func NewProcessingService(as *accounts.AccountService, ts *transfers.TransferService) *ProcessingService {
	return &ProcessingService{
		accountService:  as,
		transferService: ts,
	}
}

func (s *ProcessingService) Transaction(ctx context.Context, channel chan<- string, tr *transfer.Model) error {
	fmt.Println("started processing transaction")
	/*
		(only sender => withdraw)
		(only receiver => deposit)
		(both => transfer)
	*/
	var err error
	var senderAcc *account.Model
	var receiverAcc *account.Model

	emptyID := uuid.UUID{}

	if tr.Receiver == emptyID && tr.Sender == emptyID {
		return errors.New("cannot operate without accounts")
	}

	if tr.Receiver == tr.Sender {
		return errors.New("cannot operate between same account")
	}

	if tr.Sender != emptyID {
		senderAcc, err = s.accountService.GetByID(tr.Sender)
		if err != nil {
			return err
		}
	}

	if time.Now().After(tr.Scheduled.Time()) {
		if tr.Sender != emptyID && senderAcc.Balance-float64(tr.Amount) < 0 {
			tr, err = s.transferService.Rejecte(tr.ID)
		} else {
			tr, err = s.transferService.Approve(tr.ID)
		}
	} else {
		tr, err = s.transferService.Schedule(tr.ID)
	}

	if err != nil {
		return fmt.Errorf("tranfer failed to update status: %+v", err)
	}

	if tr.Status == transfer.Approved && tr.Receiver != emptyID {
		receiverAcc, err = s.accountService.GetByID(tr.Receiver)
		if err != nil {
			return err
		}

		receiverAcc.Balance += float64(tr.Amount)
		if err := s.accountService.Update(*receiverAcc); err != nil {
			return fmt.Errorf("couln't update sender account: %v", err)
		}
	}

	if tr.Status == transfer.Approved && tr.Sender != emptyID {
		senderAcc.Balance -= float64(tr.Amount)
		if err := s.accountService.Update(*senderAcc); err != nil {
			return fmt.Errorf("couln't update sender account: %v", err)
		}
	}

	if channel != nil {
		message, err := FormatMessage("transactionupdate", tr)
		if err != nil {
			log.Printf("fail send message: %+v", err)
		}

		channel <- message
	}

	fmt.Println("Transaction processed")
	ctx.Done()

	return nil
}

type Fail struct {
	Err      error
	Transfer *transfer.Model
}

func (s *ProcessingService) TransactionsScheduledForToday(ctx context.Context, channel chan<- string) error {
	transfers, err := s.transferService.GetTodaysPending()
	if err != nil {
		return err
	}

	var fails []Fail

	for _, transfer := range transfers {
		err := s.Transaction(ctx, channel, &transfer)
		if err != nil {
			fail := Fail{
				Err:      err,
				Transfer: &transfer,
			}

			fails = append(fails, fail)
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("failed to process %d transactions", len(fails))
	}

	return nil
}

func FormatMessage(name string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}

	buf := bytes.NewBuffer([]byte{})

	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", name))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buf.String()))

	return sb.String(), nil
}
