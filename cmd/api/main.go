package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/tsntt/fableflow/api"
	"github.com/tsntt/fableflow/data/postgres"
	"github.com/tsntt/fableflow/src/service/accounts"
	"github.com/tsntt/fableflow/src/service/banks"
	"github.com/tsntt/fableflow/src/service/processing"
	"github.com/tsntt/fableflow/src/service/transfers"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db, err := postgres.NewPostgresStore(os.Getenv("DBCONN"))
	if err != nil {
		log.Fatal(err)
	}

	transferStore := postgres.NewTransferStorage(db)
	transferService := transfers.NewTransferService(transferStore)

	accountStore := postgres.NewAccountStorage(db)
	accountService := accounts.NewAccountService(accountStore)

	bankStore := postgres.NewBankStorage(db)
	bankService := banks.NewBankService(bankStore)

	processingService := processing.NewProcessingService(accountService, transferService)

	ctx := context.Background()

	c := cron.New()
	c.AddFunc("30 8 * * 1-5", func() {
		processingService.TransactionsScheduledForToday(ctx, api.MsgChan)
	})
	c.Start()

	api := api.NewApiServer(transferService, accountService, processingService, bankService)

	// testing token things
	// tk, err := tokens.NewBankToken("www.google.com", "1234")
	// fmt.Println(tk, err)

	// tk2, err := tokens.NewAccountToken("4321", "123")
	// fmt.Println(tk2, err)

	// claims, err := tokens.VerifyToken(tk)
	// fmt.Println(claims["id"], err)

	// claims2, err := tokens.VerifyToken(tk2)
	// fmt.Println(claims2, err)
	// end token things

	api.Run()
}
