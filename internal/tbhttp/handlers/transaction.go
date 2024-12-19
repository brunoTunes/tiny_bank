package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"http/internal/service/transaction"
	"http/internal/tbhttp/handlers/request"
	"http/internal/tbhttp/handlers/response"
)

func RegisterTransactionHandler(mux *http.ServeMux, logger *slog.Logger, transactionSvc *transaction.Service) {
	logger.Debug("registering transaction endpoints")

	logger.Debug("registering POST /transaction")
	mux.Handle("POST /transaction", handlePostTransaction(logger, transactionSvc))
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
