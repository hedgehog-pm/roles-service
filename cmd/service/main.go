package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hedgehog-pm/roles-service/internal/api"
	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolerepository "github.com/hedgehog-pm/roles-service/internal/repository/role"
	assignmentrepository "github.com/hedgehog-pm/roles-service/internal/repository/roleassignment"
	templaterepository "github.com/hedgehog-pm/roles-service/internal/repository/roletemplate"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
	"github.com/hedgehog-pm/roles-service/internal/usecase/createrole"
	"github.com/hedgehog-pm/roles-service/internal/usecase/createtemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroleassignments"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroles"
	"github.com/hedgehog-pm/roles-service/internal/usecase/permission"
	"github.com/hedgehog-pm/roles-service/internal/usecase/role"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roleassignment"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roletemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/updaterole"
	"github.com/hedgehog-pm/roles-service/internal/usecase/updatetemplate"
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

	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		return err
	}

	transactor := postgres.NewTransactor(pool)
	roles := rolerepository.New(pool, catalog)
	templates := templaterepository.New(pool, catalog)
	assignments := assignmentrepository.New(pool)
	permissions := permission.New(catalog)
	roleHandler := api.NewRoleHandler(role.New(roles, transactor), listroles.New(roles, transactor, catalog), createrole.New(roles, transactor, catalog), updaterole.New(roles, transactor, catalog), permissions)
	templateHandler := api.NewRoleTemplateHandler(roletemplate.New(templates, transactor), createtemplate.New(templates, transactor, catalog), updatetemplate.New(templates, transactor, catalog), permissions)
	assignmentHandler := api.NewRoleAssignmentHandler(roleassignment.New(assignments, transactor), listroleassignments.New(assignmentrepository.NewReader(pool), transactor))
	handler := api.NewHandler(roleHandler, templateHandler, assignmentHandler, permissions)

	server, err := api.NewServer(handler)
	if err != nil {
		return err
	}

	server.Use(api.DevelopmentActor())
	slog.WarnContext(ctx, "development authentication bypass enabled", "actor_id", 1)

	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}

	return serve(ctx, &http.Server{Addr: address, Handler: server, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second})
}
