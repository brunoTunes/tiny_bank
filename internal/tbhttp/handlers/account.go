package handlers

import (
	"log/slog"
	"net/http"

	"http/internal/service/account"
	"http/internal/tbhttp/handlers/response"
)

func RegisterAccountHandler(mux *http.ServeMux, logger *slog.Logger, accountSvc *account.Service) {
	logger.Debug("registering account endpoints")

	logger.Debug("registering GET /users/{id}/accounts")
	mux.Handle("GET /users/{id}/accounts", handleGetUserAccounts(logger, accountSvc))
}

func handleGetUserAccounts(logger *slog.Logger, accountSvc *account.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.PathValue("id")
			if userID == "" {
				logger.InfoContext(r.Context(), "invalid path")
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
				return
			}

			accs, err := accountSvc.GetUserAccounts(userID)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to create user", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "failed to create user", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusOK, response.AccountsResponseFromDomain(accs))
		},
	)
}
