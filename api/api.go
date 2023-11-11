package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tsntt/fableflow/middleware"
	"github.com/tsntt/fableflow/src/service/accounts"
	"github.com/tsntt/fableflow/src/service/banks"
	"github.com/tsntt/fableflow/src/service/processing"
	"github.com/tsntt/fableflow/src/service/transfers"
	"github.com/uptrace/bunrouter"
)

var MsgChan chan string

type ApiServer struct {
	transferService      *transfers.TransferService
	accountService       *accounts.AccountService
	bankService          *banks.BankService
	processingService    *processing.ProcessingService
	trasactionsWaitGroup *sync.WaitGroup
}

func NewApiServer(ts *transfers.TransferService, as *accounts.AccountService, ps *processing.ProcessingService, bs *banks.BankService) *ApiServer {
	return &ApiServer{
		transferService:      ts,
		accountService:       as,
		processingService:    ps,
		bankService:          bs,
		trasactionsWaitGroup: &sync.WaitGroup{},
	}
}

func (srv *ApiServer) Run() {

	router := bunrouter.New(
		bunrouter.Use(middleware.Cors),
		bunrouter.Use(middleware.RateLimit),
	)

	router.POST("/register", bunrouter.HTTPHandlerFunc(srv.HandlerRequestNewBank))
	router.GET("/activate/:hash", bunrouter.HTTPHandlerFunc(srv.HandlerActivateBank))

	router.POST("/account", bunrouter.HTTPHandlerFunc(srv.HandleNewAccount))

	account := router.NewGroup("/account")
	account = account.Use(middleware.AccountAuth)
	// get on /account get logged acc
	account.GET("/:id", bunrouter.HTTPHandlerFunc(srv.HandleGetAccountByID))
	// receive update on transfers
	account.GET("/sse", bunrouter.HTTPHandlerFunc(srv.EventProcessTransactions))

	// accountAuth
	transfer := account.NewGroup("/transfer")
	// new transfer
	transfer.POST("", bunrouter.HTTPHandlerFunc(srv.HandleNewTransfer))
	// get one transfer
	transfer.GET("/:id", bunrouter.HTTPHandlerFunc(srv.HandleGetTransferByID))
	// cancel transfer if pending
	transfer.PATCH("/cancel/:id", bunrouter.HTTPHandlerFunc(srv.HandleCancelTransfer))

	transfers := account.NewGroup("/transfers")
	// account last 30 days transfers
	transfers.GET("", bunrouter.HTTPHandlerFunc(srv.HandleGetTransfersByAccount))
	// get transfer from period
	transfers.GET("/byperiod", bunrouter.HTTPHandlerFunc(srv.HandleGetTransfersByPeriod))

	srv.trasactionsWaitGroup.Wait()

	server := http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      router,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
	}

	log.Printf("Api listening at port: 4000")
	go server.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
