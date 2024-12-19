package tbhttp

import (
	"context"
	"log/slog"
	"net/http"

	"http/internal/service/account"
	"http/internal/service/transaction"
	"http/internal/service/user"
	"http/internal/tbhttp/handlers"
)

func NewServer(
	ctx context.Context,
	logger *slog.Logger,
	userService *user.Service,
	accountService *account.Service,
	transactionService *transaction.Service,
) http.Handler {
	mux := http.NewServeMux()
	handlers.RegisterUserHandler(mux, logger, userService)
	handlers.RegisterAccountHandler(mux, logger, accountService)
	handlers.RegisterTransactionHandler(mux, logger, transactionService)
	return mux
}
