package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type recordingRoles struct {
	application.RoleRepository
	saved *domain.Role
	err   error
}

func (r *recordingRoles) Insert(_ context.Context, value *domain.Role) (*domain.Role, error) {
	if r.err != nil {
		return nil, r.err
	}

	saved := *value
	saved.ID = 42
	r.saved = &saved

	return &saved, nil
}

func TestCreateValidation(t *testing.T) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeProject, 1)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name  string
		codes []domain.PermissionCode
	}{
		{"empty", nil}, {"unknown", []domain.PermissionCode{"unknown"}}, {"wrong scope", []domain.PermissionCode{domain.SystemRoleCreate}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &recordingRoles{}
			service := application.NewRoles(repo, catalog)

			_, err := service.Create(t.Context(), 1, application.CreateRoleInput{Name: "role", Scope: scope, PermissionCodes: tc.codes})
			if err == nil || repo.saved != nil {
				t.Fatal("invalid permissions reached storage")
			}
		})
	}

	repo := &recordingRoles{}
	service := application.NewRoles(repo, catalog)

	role, err := service.Create(t.Context(), 1, application.CreateRoleInput{Name: " role ", Scope: scope, PermissionCodes: []domain.PermissionCode{domain.ProjectView, domain.ProjectView}})
	if err != nil {
		t.Fatal(err)
	}

	if role.ID != 42 || role.Name != "role" || len(role.Permissions.Codes()) != 1 {
		t.Fatal("unexpected role")
	}

	repo.err = errors.New("storage unavailable")

	role, err = service.Create(t.Context(), 1, application.CreateRoleInput{Name: "role", Scope: scope, PermissionCodes: []domain.PermissionCode{domain.ProjectView}})
	if !errors.Is(err, repo.err) || role != nil {
		t.Fatal("storage failure reported as success")
	}
}
