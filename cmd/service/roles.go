package main

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/repository/postgres"
	application "github.com/hedgehog-pm/roles-service/internal/service"
)

type rolesModule struct {
	Roles       *application.Roles
	Templates   *application.Templates
	Assignments *application.Assignments
	Permissions *application.Permissions
}

func newRolesModule(pool *pgxpool.Pool) (*rolesModule, error) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		return nil, err
	}

	return &rolesModule{
		Roles:       application.NewRoles(postgres.NewRoles(pool, catalog), catalog),
		Templates:   application.NewTemplates(postgres.NewTemplates(pool, catalog), catalog),
		Assignments: application.NewAssignments(postgres.NewAssignments(pool)),
		Permissions: application.NewPermissions(catalog),
	}, nil
}
