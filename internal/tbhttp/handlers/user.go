package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"http/internal/service/user"
	"http/internal/tbhttp/handlers/request"
	"http/internal/tbhttp/handlers/response"
)

func RegisterUserHandler(mux *http.ServeMux, logger *slog.Logger, userSvc *user.Service) {
	logger.Debug("registering users endpoints")

	logger.Debug("registering POST /users")
	mux.Handle("POST /users", handlePostUsers(logger, userSvc))

	logger.Debug("registering GET /users")
	mux.Handle("GET /users", handleGetUsers(logger, userSvc))

	logger.Debug("registering DELETE /user/{id}")
	mux.Handle("DELETE /user/{id}", handleDeleteUser(logger, userSvc))
}

func handlePostUsers(logger *slog.Logger, userSvc *user.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var newUserRequest request.UserRequest

			if err := json.NewDecoder(r.Body).Decode(&newUserRequest); err != nil {
				logger.ErrorContext(r.Context(), "failed to decode json", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid json", Details: err.Error()})
				return
			}

			user, err := userSvc.CreateUser(newUserRequest.Name)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to create user", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to create user", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusCreated, response.UserResponseFromDomain(user))
		},
	)
}

func handleGetUsers(logger *slog.Logger, userSvc *user.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			returnDeletedParam := r.URL.Query().Get("return-deleted")

			var returnDeleted bool
			var err error
			if returnDeletedParam != "" {
				returnDeleted, err = strconv.ParseBool(returnDeletedParam)
				if err != nil {
					logger.InfoContext(r.Context(), "invalid return-deleted parameter", "error", err)
					writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid return-deleted parameter", Details: err.Error()})
					return
				}
			}

			usrs, err := userSvc.GetUsers(returnDeleted)
			if err != nil {
				logger.ErrorContext(r.Context(), "failed to get users", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to get users", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusOK, response.GetUsersFromDomain(usrs))
		},
	)
}

func handleDeleteUser(logger *slog.Logger, userSvc *user.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.PathValue("id")
			if userID == "" {
				logger.InfoContext(r.Context(), "invalid path")
				writeResponseJson(r.Context(), logger, w, http.StatusBadRequest, response.Error{Message: "invalid path", Details: "{id} is empty"})
				return
			}

			err := userSvc.DeleteUser(userID)
			if err != nil {
				// TODO unwrap validation error and change http status accordingly
				logger.ErrorContext(r.Context(), "failed to create user", "error", err)
				writeResponseJson(r.Context(), logger, w, http.StatusUnprocessableEntity, response.Error{Message: "failed to create user", Details: err.Error()})
				return
			}

			writeResponseJson(r.Context(), logger, w, http.StatusNoContent, nil)
		},
	)
}
