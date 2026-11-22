package postgres_test

import (
	"errors"
	"testing"

	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
	"github.com/hedgehog-pm/roles-service/internal/roles/infrastructure/postgres"
)

func TestClientPagesAndAssignments(t *testing.T) {
	pool := testDatabase(t)

	module, err := testModule(pool)
	if err != nil {
		t.Fatal(err)
	}

	scope, err := domain.NewScopeRef(domain.ScopeProject, testID(t, pool))
	if err != nil {
		t.Fatal(err)
	}

	actor, user := testID(t, pool), testID(t, pool)

	var created []*domain.Role

	for _, name := range []string{"Same", "Same", "literal%_name"} {
		v, err := module.Roles.Create(t.Context(), actor, application.CreateRoleInput{Name: name, Scope: scope, PermissionCodes: []domain.PermissionCode{domain.ProjectView}})
		if err != nil {
			t.Fatal(err)
		}

		created = append(created, v)
	}

	options := application.ListOptions{Limit: 1, Sort: "name", Direction: "asc"}

	newest, err := module.Roles.Page(t.Context(), scope, application.ListOptions{Limit: 1})
	if err != nil || len(newest.Items) != 1 || newest.Items[0].ID != created[2].ID {
		t.Fatalf("default order: %+v %v", newest, err)
	}

	page, err := module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Total != 3 || len(page.Items) != 1 || page.Items[0].ID != created[0].ID {
		t.Fatalf("first page: %+v %v", page, err)
	}

	options.Offset = 1

	page, err = module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Items[0].ID != created[1].ID {
		t.Fatalf("tie order: %+v %v", page, err)
	}

	options.Offset = 100

	page, err = module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Total != 3 || len(page.Items) != 0 {
		t.Fatalf("empty page: %+v %v", page, err)
	}

	options.Offset = 0
	options.Search = "%_"
	options.PermissionCode = domain.ProjectView

	page, err = module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Total != 1 || page.Items[0].ID != created[2].ID {
		t.Fatalf("literal search: %+v %v", page, err)
	}

	options.Search = "sAmE"

	page, err = module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Total != 2 {
		t.Fatalf("case insensitive search: %+v %v", page, err)
	}

	options.PermissionCode = domain.ProjectEdit

	page, err = module.Roles.Page(t.Context(), scope, options)
	if err != nil || page.Total != 0 {
		t.Fatalf("permission filter: %+v %v", page, err)
	}

	type outcome struct {
		result application.AssignmentResult
		err    error
	}

	outcomes := make(chan outcome, 8)

	for range 8 {
		go func() {
			result, err := module.Assignments.Ensure(t.Context(), actor, created[0].ID, user)
			outcomes <- outcome{result, err}
		}()
	}

	count := 0

	var original domain.RoleAssignment

	for range 8 {
		v := <-outcomes
		if v.err != nil {
			t.Fatal(v.err)
		}

		if v.result.Created {
			count++
			original = v.result.Assignment
		}
	}

	if count != 1 {
		t.Fatalf("created %d assignments", count)
	}

	again, err := module.Assignments.Ensure(t.Context(), actor+1, created[0].ID, user)
	if err != nil || again.Created || again.Assignment.AssignedBy() != actor || !again.Assignment.AssignedAt().Equal(original.AssignedAt()) {
		t.Fatalf("idempotency: %+v %v", again, err)
	}

	var events int
	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM outbox_messages WHERE aggregate_id=$1 AND event_type='RoleAssigned'", created[0].ID).Scan(&events); err != nil || events != 1 {
		t.Fatalf("events %d: %v", events, err)
	}

	assignments, err := module.Assignments.Page(t.Context(), application.AssignmentFilter{UserID: user, ScopeType: scope.Type(), ScopeID: scope.ID(), Limit: 20})
	if err != nil || assignments.Total != 1 || assignments.Items[0].RoleName != "Same" {
		t.Fatalf("user assignments: %+v %v", assignments, err)
	}

	empty, err := module.Assignments.Page(t.Context(), application.AssignmentFilter{UserID: user, ScopeType: domain.ScopeWorkspace})
	if err != nil || empty.Total != 0 {
		t.Fatalf("scope filter: %+v %v", empty, err)
	}

	for range 2 {
		if err := module.Assignments.Revoke(t.Context(), actor, created[0].ID, user); err != nil {
			t.Fatal(err)
		}
	}

	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM outbox_messages WHERE aggregate_id=$1 AND event_type='RoleRevoked'", created[0].ID).Scan(&events); err != nil || events != 1 {
		t.Fatalf("revoke events %d: %v", events, err)
	}

	system := testID(t, pool)

	template, err := module.Templates.Create(t.Context(), actor, application.CreateTemplateInput{SystemID: system, Name: "Template_%", TargetScopeType: domain.ScopeProject, PermissionCodes: []domain.PermissionCode{domain.ProjectView}})
	if err != nil {
		t.Fatal(err)
	}

	templates, err := module.Templates.Page(t.Context(), application.TemplateFilter{SystemID: system, TargetScopeType: domain.ScopeProject}, application.ListOptions{Search: "_%", Limit: 1, Offset: 10})
	if err != nil || templates.Total != 1 || len(templates.Items) != 0 {
		t.Fatalf("template page: %+v %v", templates, err)
	}

	if err := module.Templates.Delete(t.Context(), actor, template.ID, template.Version); err != nil {
		t.Fatal(err)
	}

	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	if err := module.Roles.Delete(t.Context(), actor, created[0].ID, created[0].Version); err != nil {
		t.Fatal(err)
	}

	if err := postgres.NewRoles(pool, catalog).Update(t.Context(), created[0], created[0].Version); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("update after concurrent deletion: %v", err)
	}

	if err := module.Assignments.Revoke(t.Context(), actor, created[0].ID, user); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("revoke for missing role: %v", err)
	}
}
