package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tsntt/fableflow/src/service/banks"
	"github.com/tsntt/fableflow/src/util"
)

func (s *ApiServer) HandlerRequestNewBank(w http.ResponseWriter, r *http.Request) {
	req := banks.NewBankReqDTO{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := s.bankService.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link := `http://localhost:4000/activate/` + res.Hash
	go util.SendSimpleMessage(link)

	util.WriteJson(w, http.StatusCreated, "")
}

func (s *ApiServer) HandlerActivateBank(w http.ResponseWriter, r *http.Request) {
	hash := strings.TrimPrefix(r.URL.Path, "/activate/")

	err := s.bankService.Activate(hash)
	if err != nil {
		http.Error(w, "couldn't activate bank", http.StatusInternalServerError)
		return
	}

	util.WriteJson(w, http.StatusOK, "Activated")
}
