package roles

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
	"github.com/hedgehog-pm/roles-service/internal/roles/infrastructure/postgres"
)

type Module struct {
	Roles       *application.Roles
	Templates   *application.Templates
	Assignments *application.Assignments
	Permissions *application.Permissions
}

func New(pool *pgxpool.Pool) (*Module, error) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		return nil, err
	}

	return &Module{
		Roles:       application.NewRoles(postgres.NewRoles(pool, catalog), catalog),
		Templates:   application.NewTemplates(postgres.NewTemplates(pool, catalog), catalog),
		Assignments: application.NewAssignments(postgres.NewAssignments(pool)),
		Permissions: application.NewPermissions(catalog),
	}, nil
}
