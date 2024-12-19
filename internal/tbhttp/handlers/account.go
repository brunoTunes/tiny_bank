package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"http/internal/service/account"
	"http/internal/service/transaction"
	"http/internal/tbhttp/handlers/request"
	"http/internal/tbhttp/handlers/response"
)

func RegisterAccountHandler(mux *http.ServeMux, logger *slog.Logger, accountSvc *account.Service, transactionSvc *transaction.Service) {
	logger.Debug("registering account endpoints")

	logger.Debug("registering POST /account/{id}/balance")
	mux.Handle("POST /account/{id}/balance", handlePostAccountBalance(logger, accountSvc))

	logger.Debug("registering GET /account/{id}/transactions")
	mux.Handle("GET /account/{id}/transactions", handleGetAccountTransactions(logger, transactionSvc))
}

func handlePostAccountBalance(logger *slog.Logger, accountSvc *account.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var postAccountBalance request.AccountBalance

			accountID := r.PathValue("id")
			if accountID == "" {
				logger.InfoContext(r.Context(), "invalid path")
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
				return
			}

			if err := json.NewDecoder(r.Body).Decode(&postAccountBalance); err != nil {
				logger.ErrorContext(r.Context(), "failed to decode json", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid json", Details: err.Error()})
				return
			}

			account, err := accountSvc.AddBalance(accountID, postAccountBalance.Balance)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to add balance to account", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to add balance to account", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusCreated, response.AccountResponseFromDomain(account))
		},
	)
}

func handleGetAccountTransactions(logger *slog.Logger, transactionSvc *transaction.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fromDateParam := r.URL.Query().Get("from-date")
		toDateParam := r.URL.Query().Get("to-date")

		// TODO generalize query params parsing
		now := time.Now()
		var fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		var err error
		if fromDateParam != "" {
			fromDate, err = time.Parse("2006-01-02", fromDateParam)
			if err != nil {
				logger.InfoContext(r.Context(), "invalid from date parameter", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid fromDate parameter", Details: err.Error()})
				return
			}
		}

		var toDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
		if toDateParam != "" {
			toDate, err = time.Parse("2006-01-02", toDateParam)
			if err != nil {
				logger.InfoContext(r.Context(), "invalid to date parameter", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid to date parameter", Details: err.Error()})
				return
			}
		}

		accountID := r.PathValue("id")
		if accountID == "" {
			logger.InfoContext(r.Context(), "invalid path")
			writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
			return
		}

		transactions, err := transactionSvc.GetAccountTransactionHistory(accountID, fromDate, toDate)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to get account transactions", "error", err)
			writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "failed to get account transaction history", Details: err.Error()})
			return
		}

		writeResponseJson(r.Context(), logger, w, http.StatusOK, response.AccountTransactionsHistoryFromDomain(transactions))
	})
}
