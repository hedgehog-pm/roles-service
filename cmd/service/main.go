package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hedgehog-pm/roles-service/internal/api"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startup, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := postgres.Open(startup)
	if err != nil {
		return err
	}
	defer pool.Close()

	module, err := newRolesModule(pool)
	if err != nil {
		return err
	}

	var authenticate api.Authenticate

	if rawID := os.Getenv("DEV_ACTOR_ID"); rawID != "" {
		id, err := strconv.Atoi(rawID)
		if err != nil || id <= 0 || id > 2147483647 {
			return errors.New("DEV_ACTOR_ID must be an integer between 1 and 2147483647")
		}

		authenticate = func(context.Context, *http.Request) (int, error) { return id, nil }

		slog.WarnContext(ctx, "local development authentication bypass enabled", "actor_id", id)
	}

	server, err := api.NewServer(api.NewHandler(module.Roles, module.Templates, module.Assignments, module.Permissions), authenticate)
	if err != nil {
		return err
	}

	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}

	return serve(ctx, &http.Server{Addr: address, Handler: server, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second})
}
