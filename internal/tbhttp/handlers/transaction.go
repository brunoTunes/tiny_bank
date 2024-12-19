package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"http/internal/service/transaction"
	"http/internal/tbhttp/handlers/request"
	"http/internal/tbhttp/handlers/response"
)

func RegisterTransactionHandler(mux *http.ServeMux, logger *slog.Logger, transactionSvc *transaction.Service) {
	logger.Debug("registering transaction endpoints")

	logger.Debug("registering POST /transaction")
	mux.Handle("POST /transaction", handlePostTransaction(logger, transactionSvc))

	logger.Debug("registering POST /account/{id}/withdraw")
	mux.Handle("POST /account/{id}/withdraw", handlePostWithdraw(logger, transactionSvc))

	logger.Debug("registering POST /account/{id}/withdraw")
	mux.Handle("POST /account/{id}/deposit", handlePostDeposit(logger, transactionSvc))

	logger.Debug("registering GET /account/{id}/transactions")
	mux.Handle("GET /account/{id}/transactions", handleGetAccountTransactions(logger, transactionSvc))
}

func handlePostWithdraw(logger *slog.Logger, transactionSvc *transaction.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var postWithdraw request.Withdraw

			accountID := r.PathValue("id")
			if accountID == "" {
				logger.InfoContext(r.Context(), "invalid path")
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
				return
			}

			if err := json.NewDecoder(r.Body).Decode(&postWithdraw); err != nil {
				logger.ErrorContext(r.Context(), "failed to decode json", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid json", Details: err.Error()})
				return
			}

			tr, err := transactionSvc.Withdraw(accountID, postWithdraw.Amount)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to perform transfer", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to perform transfer", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusCreated, response.WithdrawFromDomain(tr))
		},
	)
}

func handlePostDeposit(logger *slog.Logger, transactionSvc *transaction.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var postDeposit request.Deposit

			accountID := r.PathValue("id")
			if accountID == "" {
				logger.InfoContext(r.Context(), "invalid path")
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
				return
			}

			if err := json.NewDecoder(r.Body).Decode(&postDeposit); err != nil {
				logger.ErrorContext(r.Context(), "failed to decode json", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid json", Details: err.Error()})
				return
			}

			tr, err := transactionSvc.Deposit(accountID, postDeposit.Amount)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to perform transfer", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to perform transfer", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusCreated, response.DepositFromDomain(tr))
		},
	)
}

func handlePostTransaction(logger *slog.Logger, transactionSvc *transaction.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var postTransaction request.Transaction

			if err := json.NewDecoder(r.Body).Decode(&postTransaction); err != nil {
				logger.ErrorContext(r.Context(), "failed to decode json", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid json", Details: err.Error()})
				return
			}

			tr, err := transactionSvc.Transfer(postTransaction.FromAccount, postTransaction.ToAccount, postTransaction.Amount)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to perform transfer", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to perform transfer", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusCreated, response.TransactionFromDomain(tr))
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

		writeResponseJson(r.Context(), logger, w, http.StatusOK, response.TransactionsHistoryFromDomain(transactions))
	})
}
