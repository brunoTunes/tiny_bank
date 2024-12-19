package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"http/internal/repository/memory"
	"http/internal/service/account"
	"http/internal/service/transaction"
	"http/internal/service/user"
	"http/internal/tbhttp"

	"github.com/Netflix/go-env"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL,default=debug"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		return err
	}

	logger := initLogger(config.LogLevel)

	userRepo := memory.NewUserRepository()
	accountRepo := memory.NewAccountRepository()
	transactionRepo := memory.NewTransactionRepository()
	accountService := account.NewService(accountRepo)
	userSvc := user.NewService(userRepo, accountService)
	transactionSvc := transaction.NewService(accountService, transactionRepo)

	server := &http.Server{
		Addr:    ":8080",
		Handler: tbhttp.NewServer(ctx, logger, userSvc, accountService, transactionSvc),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "HTTP server error", "error", err)
		}

		logger.InfoContext(ctx, "Stopped serving new connections.")
	}()

	<-ctx.Done()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.ErrorContext(ctx, "HTTP shutdown error", "error", err)
	}

	logger.InfoContext(ctx, "Graceful shutdown complete.")

	return nil
}

func initLogger(logLevelEnv string) *slog.Logger {
	logLevel := slog.LevelDebug

	switch logLevelEnv {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
