package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// missing fallback, log error for now
func writeResponseJson(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
