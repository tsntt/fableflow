package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/middleware"
	"github.com/tsntt/fableflow/src/domain/transfer"
	"github.com/tsntt/fableflow/src/service/transfers"
	"github.com/tsntt/fableflow/src/util"
)

type FullTransferResDTO struct {
	ID        uuid.UUID `json:"id"`
	Receiver  uuid.UUID `json:"receiver"`
	Sender    uuid.UUID `json:"sender"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Scheduled time.Time `json:"scheduled"`
	CreatedAt time.Time `json:"created_at"`
}

func (srv *ApiServer) HandleNewTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.WriteJson(w, http.StatusOK, "")
		return
	}
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newTransfer transfers.NewTransferReqDTO

	err := json.NewDecoder(r.Body).Decode(&newTransfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accID := uuid.MustParse(r.Context().Value(middleware.AccID).(string))

	if (newTransfer.Sender != uuid.UUID{}) && newTransfer.Sender != accID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tr, err := srv.transferService.CreateTransfer(newTransfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	srv.trasactionsWaitGroup.Add(1)
	go func() {
		defer srv.trasactionsWaitGroup.Done()

		ctx := context.Background()
		err := srv.processingService.Transaction(ctx, MsgChan, tr)
		if err != nil {
			log.Println(err)
		}
	}()

	res := transfers.NewTransferRespDTO{
		ID:     tr.ID,
		Status: string(tr.Status),
	}

	util.WriteJson(w, http.StatusOK, res)
}

func (srv *ApiServer) HandleCancelTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.WriteJson(w, http.StatusOK, "")
		return
	}
	if r.Method != "PATCH" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/account/tranfer/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := srv.transferService.Cancel(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := transfers.NewTransferRespDTO{
		ID:     id,
		Status: string(transfer.Canceled),
	}

	util.WriteJson(w, http.StatusOK, res)
}

func (srv *ApiServer) HandleGetTransferByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/account/transfer/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transfer, err := srv.transferService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accID := uuid.MustParse(r.Context().Value(middleware.AccID).(string))

	if transfer.Sender != accID && transfer.Receiver != accID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res := FullTransferResDTO{
		ID:        transfer.ID,
		Sender:    transfer.Sender,
		Receiver:  transfer.Receiver,
		Amount:    float64(transfer.Amount),
		Status:    string(transfer.Status),
		Scheduled: transfer.Scheduled.Time(),
		CreatedAt: transfer.CreatedAt,
	}

	util.WriteJson(w, http.StatusOK, res)
}
func (srv *ApiServer) HandleGetTransfersByAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accID := uuid.MustParse(r.Context().Value(middleware.AccID).(string))

	ts, err := srv.transferService.GetByAccountAndPeriod(accID, time.Now().AddDate(0, -1, 0), time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]FullTransferResDTO, len(ts))

	for i, transfer := range ts {
		res[i] = FullTransferResDTO{
			ID:        transfer.ID,
			Sender:    transfer.Sender,
			Receiver:  transfer.Receiver,
			Amount:    float64(transfer.Amount),
			Status:    string(transfer.Status),
			Scheduled: transfer.Scheduled.Time(),
			CreatedAt: transfer.CreatedAt,
		}
	}

	util.WriteJson(w, http.StatusOK, res)
}
func (srv *ApiServer) HandleGetTransfersByPeriod(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accID := uuid.MustParse(r.Context().Value(middleware.AccID).(string))

	var dates transfer.Period
	if err := json.NewDecoder(r.Body).Decode(&dates); err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad formated date. Expect string format: 2023-11-01T18:30:00Z", http.StatusBadRequest)
		return
	}

	ts, err := srv.transferService.GetByAccountAndPeriod(accID, dates.Start, dates.End)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]FullTransferResDTO, len(ts))

	for i, transfer := range ts {
		res[i] = FullTransferResDTO{
			ID:        transfer.ID,
			Sender:    transfer.Sender,
			Receiver:  transfer.Receiver,
			Amount:    float64(transfer.Amount),
			Status:    string(transfer.Status),
			Scheduled: transfer.Scheduled.Time(),
			CreatedAt: transfer.CreatedAt,
		}
	}

	util.WriteJson(w, http.StatusOK, res)
}
