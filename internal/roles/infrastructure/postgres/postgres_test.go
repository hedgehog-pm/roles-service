package postgres_test

import (
	"context"
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
	"github.com/hedgehog-pm/roles-service/internal/roles/infrastructure/postgres"
)

type moduleFixture struct {
	Roles       *application.Roles
	Templates   *application.Templates
	Assignments *application.Assignments
}

func testModule(pool *pgxpool.Pool) (*moduleFixture, error) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		return nil, err
	}

	return &moduleFixture{Roles: application.NewRoles(postgres.NewRoles(pool, catalog), catalog), Templates: application.NewTemplates(postgres.NewTemplates(pool, catalog), catalog), Assignments: application.NewAssignments(postgres.NewAssignments(pool))}, nil
}

func testDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("ROLES_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("ROLES_TEST_DATABASE_URL is required for PostgreSQL integration tests")
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(config.Database, "_test") {
		t.Fatal("test database name must end in _test")
	}

	db := stdlib.OpenDB(*config)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("../../../../migrations/postgresql/master"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err = provider.Up(t.Context()); err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(t.Context(), dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(pool.Close)

	if _, err := pool.Exec(t.Context(), "CREATE SEQUENCE IF NOT EXISTS test_reference_ids AS integer"); err != nil {
		t.Fatal(err)
	}

	return pool
}

func TestRoleModulePostgres(t *testing.T) {
	pool := testDatabase(t)

	module, err := testModule(pool)
	if err != nil {
		t.Fatal(err)
	}

	actor := testID(t, pool)

	scope, err := domain.NewScopeRef(domain.ScopeProject, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	ctx := t.Context()

	role, err := module.Roles.Create(ctx, actor, application.CreateRoleInput{Name: "developer", Scope: scope, PermissionCodes: []domain.PermissionCode{domain.ProjectView}})
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := module.Roles.Get(ctx, role.ID)
	if err != nil {
		t.Fatal(err)
	}

	if role.ID <= 0 || loaded.ID != role.ID || loaded.Version != role.Version || loaded.CreatedAt.UnixMicro() != role.CreatedAt.UnixMicro() || loaded.UpdatedAt.UnixMicro() != role.UpdatedAt.UnixMicro() || !loaded.Permissions.Equal(role.Permissions) {
		t.Fatal("round trip changed persisted state")
	}

	list, err := module.Roles.Page(ctx, scope, application.ListOptions{})
	if err != nil || len(list.Items) != 1 {
		t.Fatalf("list: %v, %d", err, len(list.Items))
	}

	input := application.UpdateRoleInput{ID: role.ID, Version: role.Version, Name: "renamed", PermissionCodes: []domain.PermissionCode{domain.ProjectEdit}}

	updated, err := module.Roles.Update(ctx, actor, input)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Version != 3 {
		t.Fatalf("expected both domain mutations, got version %d", updated.Version)
	}

	if _, err = module.Roles.Update(ctx, actor, input); !errors.Is(err, application.ErrConflict) {
		t.Fatalf("stale update: %v", err)
	}

	input.Version = updated.Version

	same, err := module.Roles.Update(ctx, actor, input)
	if err != nil || same.Version != updated.Version {
		t.Fatalf("no-op update: %v", err)
	}

	var count int
	if err = pool.QueryRow(ctx, "SELECT count(*) FROM outbox_messages WHERE aggregate_id=$1", role.ID).Scan(&count); err != nil || count != 2 {
		t.Fatalf("no-op produced event: %d %v", count, err)
	}

	user := testID(t, pool)
	if err = module.Assignments.Assign(ctx, actor, role.ID, user); err != nil {
		t.Fatal(err)
	}

	if err = module.Assignments.Assign(ctx, actor, role.ID, user); err != nil {
		t.Fatalf("duplicate assignment: %v", err)
	}

	assigned, err := module.Assignments.List(ctx, role.ID)
	if err != nil || len(assigned) != 1 || assigned[0].AssignedBy() != actor {
		t.Fatalf("assignments: %v", err)
	}

	if err = module.Assignments.Revoke(ctx, actor, role.ID, user); err != nil {
		t.Fatal(err)
	}

	if err = module.Assignments.Assign(ctx, actor, role.ID, user); err != nil {
		t.Fatal(err)
	}

	if err = module.Roles.Delete(ctx, actor, role.ID, updated.Version); err != nil {
		t.Fatal(err)
	}

	if _, err = module.Roles.Get(ctx, role.ID); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("deleted role: %v", err)
	}

	if err = pool.QueryRow(ctx, "SELECT count(*) FROM role_assignments WHERE role_id=$1", role.ID).Scan(&count); err != nil || count != 0 {
		t.Fatalf("assignment cascade: %v", err)
	}

	if err = pool.QueryRow(ctx, "SELECT count(*) FROM role_permissions WHERE role_id=$1", role.ID).Scan(&count); err != nil || count != 0 {
		t.Fatalf("permission cascade: %v", err)
	}
}

func TestDatabaseGeneratedIdentity(t *testing.T) {
	pool := testDatabase(t)

	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeProject, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	permissions, err := catalog.NewSet(domain.ScopeProject, []domain.PermissionCode{domain.ProjectView})
	if err != nil {
		t.Fatal(err)
	}

	value, err := domain.NewRole("generated ID", scope, permissions, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	repository := postgres.NewRoles(pool, catalog)

	first, err := repository.Insert(t.Context(), value)
	if err != nil {
		t.Fatal(err)
	}

	second, err := repository.Insert(t.Context(), value)
	if err != nil {
		t.Fatal(err)
	}

	if value.ID != 0 || first.ID <= 0 || second.ID <= 0 || first.ID == second.ID {
		t.Fatal("inserts must return distinct database IDs without changing the input")
	}
}

func TestTemplatesPostgres(t *testing.T) {
	pool := testDatabase(t)

	module, err := testModule(pool)
	if err != nil {
		t.Fatal(err)
	}

	actor, system := testID(t, pool), testID(t, pool)

	template, err := module.Templates.Create(t.Context(), actor, application.CreateTemplateInput{Name: "basic", SystemID: system, TargetScopeType: domain.ScopeWorkspace, PermissionCodes: []domain.PermissionCode{domain.WorkspaceView}})
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := module.Templates.Get(t.Context(), template.ID)
	if err != nil {
		t.Fatal(err)
	}

	if template.ID <= 0 || loaded.ID != template.ID || loaded.Version != 1 || loaded.CreatedAt.UnixMicro() != template.CreatedAt.UnixMicro() {
		t.Fatal("template state changed")
	}

	list, err := module.Templates.Page(t.Context(), application.TemplateFilter{SystemID: system, TargetScopeType: domain.ScopeWorkspace}, application.ListOptions{})
	if err != nil || len(list.Items) != 1 {
		t.Fatalf("templates: %v", err)
	}

	updated, err := module.Templates.Update(t.Context(), actor, application.UpdateTemplateInput{ID: template.ID, Version: 1, Name: "edited", PermissionCodes: []domain.PermissionCode{domain.WorkspaceEdit}})
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeWorkspace, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	role, err := module.Roles.Create(t.Context(), actor, application.CreateRoleInput{Name: "from form", Scope: scope, PermissionCodes: updated.Permissions.Codes()})
	if err != nil {
		t.Fatal(err)
	}

	if err = module.Templates.Delete(t.Context(), actor, template.ID, updated.Version); err != nil {
		t.Fatal(err)
	}

	if _, err = module.Roles.Get(t.Context(), role.ID); err != nil {
		t.Fatal("template deletion affected role", err)
	}
}

func TestTransactionRollback(t *testing.T) {
	pool := testDatabase(t)

	module, err := testModule(pool)
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeProject, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}
	// A failing outbox constraint must roll back the preceding role and permission inserts.
	_, err = pool.Exec(t.Context(), "ALTER TABLE outbox_messages ADD CONSTRAINT reject_test_scope CHECK (payload->>'scope_id' <> '"+strconv.Itoa(scope.ID())+"')")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(context.WithoutCancel(t.Context()), "ALTER TABLE outbox_messages DROP CONSTRAINT reject_test_scope"); err != nil {
			t.Error(err)
		}
	})

	_, err = module.Roles.Create(t.Context(), testID(t, pool), application.CreateRoleInput{Name: "rollback", Scope: scope, PermissionCodes: []domain.PermissionCode{domain.ProjectView}})
	if err == nil {
		t.Fatal("expected outbox failure")
	}

	var count int
	if err = pool.QueryRow(t.Context(), "SELECT count(*) FROM roles WHERE scope_id=$1", scope.ID()).Scan(&count); err != nil || count != 0 {
		t.Fatalf("role persisted after failed event: %d %v", count, err)
	}
}

func TestMutationRollbackAndOptimisticLock(t *testing.T) {
	pool := testDatabase(t)

	module, err := testModule(pool)
	if err != nil {
		t.Fatal(err)
	}

	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	repository := postgres.NewRoles(pool, catalog)
	actor, user := testID(t, pool), testID(t, pool)

	scope, err := domain.NewScopeRef(domain.ScopeSystem, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	role, err := module.Roles.Create(t.Context(), actor, application.CreateRoleInput{Name: "original", Scope: scope, PermissionCodes: []domain.PermissionCode{domain.SystemRoleView}})
	if err != nil {
		t.Fatal(err)
	}

	first, err := repository.Get(t.Context(), role.ID)
	if err != nil {
		t.Fatal(err)
	}

	stale, err := repository.Get(t.Context(), role.ID)
	if err != nil {
		t.Fatal(err)
	}

	if err = first.Rename("first"); err != nil {
		t.Fatal(err)
	}

	if err = stale.Rename("stale"); err != nil {
		t.Fatal(err)
	}

	if err = repository.Update(t.Context(), first, 1); err != nil {
		t.Fatal(err)
	}

	if err = repository.Update(t.Context(), stale, 1); !errors.Is(err, application.ErrConflict) {
		t.Fatalf("repository lost update: %v", err)
	}

	if err = module.Assignments.Assign(t.Context(), actor, role.ID, user); err != nil {
		t.Fatal(err)
	}

	_, err = pool.Exec(t.Context(), "ALTER TABLE outbox_messages ADD CONSTRAINT reject_test_aggregate CHECK (aggregate_id <> '"+strconv.Itoa(role.ID)+"'::int) NOT VALID")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(context.WithoutCancel(t.Context()), "ALTER TABLE outbox_messages DROP CONSTRAINT reject_test_aggregate"); err != nil {
			t.Error(err)
		}
	})

	operations := []struct {
		name string
		run  func() error
	}{
		{"update", func() error {
			_, err := module.Roles.Update(t.Context(), actor, application.UpdateRoleInput{ID: role.ID, Version: 2, Name: "rollback", PermissionCodes: []domain.PermissionCode{domain.SystemRoleCreate}})
			return err
		}},
		{"assign", func() error { return module.Assignments.Assign(t.Context(), actor, role.ID, testID(t, pool)) }},
		{"revoke", func() error { return module.Assignments.Revoke(t.Context(), actor, role.ID, user) }},
		{"delete", func() error { return module.Roles.Delete(t.Context(), actor, role.ID, 2) }},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			if err := operation.run(); err == nil {
				t.Fatal("expected event persistence failure")
			}
		})
	}

	loaded, err := repository.Get(t.Context(), role.ID)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Name != "first" || loaded.Version != 2 || !slices.Contains(loaded.Permissions.Codes(), domain.SystemRoleView) {
		t.Fatal("failed operation changed role")
	}

	assignments, err := module.Assignments.List(t.Context(), role.ID)
	if err != nil || len(assignments) != 1 || assignments[0].UserID() != user {
		t.Fatalf("failed operation changed assignments: %v", err)
	}
}

func testID(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	var id int
	if err := pool.QueryRow(t.Context(), "SELECT nextval('test_reference_ids')").Scan(&id); err != nil {
		t.Fatal(err)
	}

	return id
}
