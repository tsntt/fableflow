package api

import (
	"fmt"
	"net/http"
)

func (srv *ApiServer) EventProcessTransactions(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	MsgChan = make(chan string)

	defer func() {
		close(MsgChan)
		fmt.Println("Client closed connection")
	}()

	for {
		select {
		case message := <-MsgChan:
			fmt.Fprintf(w, "%v\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}
}
