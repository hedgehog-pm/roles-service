package application_test

import (
	"context"
	"testing"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type updateRoles struct {
	application.RoleRepository
	value *domain.Role
}

func (r *updateRoles) Get(context.Context, int) (*domain.Role, error) { return r.value, nil }
func (r *updateRoles) Update(_ context.Context, value *domain.Role, expected int64) error {
	r.value = value

	if expected != 7 {
		panic("unexpected version")
	}

	return nil
}

type updateTemplates struct {
	application.TemplateRepository
	value *domain.RoleTemplate
}

func (r *updateTemplates) Get(context.Context, int) (*domain.RoleTemplate, error) {
	return r.value, nil
}

func (r *updateTemplates) Update(_ context.Context, value *domain.RoleTemplate, expected int64) error {
	r.value = value

	if expected != 7 {
		panic("unexpected version")
	}

	return nil
}

func TestUpdateVersion(t *testing.T) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeProject, 1)
	if err != nil {
		t.Fatal(err)
	}

	original, err := catalog.NewSet(domain.ScopeProject, []domain.PermissionCode{domain.ProjectView})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name                string
		rename, permissions bool
		version             int64
	}{
		{"unchanged", false, false, 7}, {"name", true, false, 8}, {"permissions", false, true, 8}, {"both", true, true, 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "original"
			codes := []domain.PermissionCode{domain.ProjectView}

			if tc.rename {
				name = "renamed"
			}

			if tc.permissions {
				codes = append(codes, domain.ProjectRoleCreate)
			}

			roles := &updateRoles{value: &domain.Role{ID: 1, Name: "original", Scope: scope, Permissions: original, Version: 7}}

			role, err := application.NewRoles(roles, catalog).Update(t.Context(), 1, application.UpdateRoleInput{ID: 1, Version: 7, Name: name, PermissionCodes: codes})
			if err != nil {
				t.Fatal(err)
			}

			if role.Version != tc.version {
				t.Fatalf("role version = %d, want %d", role.Version, tc.version)
			}

			templates := &updateTemplates{value: &domain.RoleTemplate{ID: 1, Name: "original", TargetScopeType: domain.ScopeProject, Permissions: original, Version: 7}}

			template, err := application.NewTemplates(templates, catalog).Update(t.Context(), 1, application.UpdateTemplateInput{ID: 1, Version: 7, Name: name, PermissionCodes: codes})
			if err != nil {
				t.Fatal(err)
			}

			if template.Version != tc.version {
				t.Fatalf("template version = %d, want %d", template.Version, tc.version)
			}
		})
	}
}
