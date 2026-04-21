package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mitayush/visa-txn/internal/config"
	"github.com/mitayush/visa-txn/internal/handler"
	"github.com/mitayush/visa-txn/internal/repository"
	"github.com/mitayush/visa-txn/internal/service"
	"github.com/mitayush/visa-txn/migrations"
)

func main() {
	r := mux.NewRouter()
	cfg := config.Load()
	registerRoutes(r, cfg)

	log.Println("Server running on", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func registerRoutes(r *mux.Router, cfg *config.Config) {
	db := repository.NewDB(cfg.DBUrl)
	if _, err := db.Exec(migrations.Init); err != nil {
		log.Fatal(err)
	}

	accountRepo := repository.NewAccountRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	txAuditRepo := repository.NewTransactionAuditRepository(db)
	accountService := service.NewAccountService(accountRepo)
	txService := service.NewTransactionService(db, txRepo, accountRepo, txAuditRepo, cfg)

	accountHandler := handler.NewAccountHandler(accountService)
	txHandler := handler.NewTransactionHandler(txService)

	r.HandleFunc("/accounts", accountHandler.Create).Methods("POST")
	r.HandleFunc("/accounts/{id}", accountHandler.Get).Methods("GET")
	r.HandleFunc("/transactions", txHandler.CreateTransaction).Methods("POST")
}
