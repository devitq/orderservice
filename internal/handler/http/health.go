package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewHealthHandler(db *sqlx.DB, redisDB *redis.Client) *HealthHandler {
	return &HealthHandler{
		DB:    db,
		Redis: redisDB,
	}
}

type HealthResponse struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	details := map[string]string{}

	if err := h.DB.PingContext(ctx); err != nil {
		details["postgres"] = "unhealthy: " + err.Error()
	} else {
		details["postgres"] = "ok"
	}

	if err := h.Redis.Ping(ctx).Err(); err != nil {
		details["redis"] = "unhealthy: " + err.Error()
	} else {
		details["redis"] = "ok"
	}

	status := "ok"
	for _, v := range details {
		if v != "ok" {
			status = "unhealthy"
			break
		}
	}

	resp := HealthResponse{Status: status, Details: details}

	w.Header().Set("Content-Type", "application/json")
	if status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(resp)
}
