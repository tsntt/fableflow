package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tsntt/fableflow/middleware"
	"github.com/tsntt/fableflow/src/service/accounts"
	"github.com/tsntt/fableflow/src/service/transfers"
	"github.com/tsntt/fableflow/src/util"
)

func (srv *ApiServer) HandleNewAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.WriteJson(w, http.StatusOK, "")
		return
	}
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newAccount accounts.NewAccountReqDTO

	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bk, err := srv.bankService.Get(r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	newAccount.BankID = bk.ID

	res, err := srv.accountService.Create(newAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newTransfer := transfers.NewTransferReqDTO{
		Receiver: res.ID,
		Amount:   newAccount.InitialDeposit,
	}
	tr, err := srv.transferService.CreateTransfer(newTransfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	token, err := util.NewAccountToken(bk.ID.String(), res.ID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Fableflowaid", token)
	util.WriteJson(w, http.StatusCreated, res)
}

func (srv *ApiServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/account/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accID := uuid.MustParse(r.Context().Value(middleware.AccID).(string))

	if id != accID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	acc, err := srv.accountService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	util.WriteJson(w, http.StatusOK, acc)
}
