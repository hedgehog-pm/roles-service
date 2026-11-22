package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hedgehog-pm/roles-service/internal/api"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type roleStore struct {
	application.RoleRepository
	saved    *domain.Role
	failure  error
	expected int64
}

func (s *roleStore) Insert(_ context.Context, v *domain.Role) (*domain.Role, error) {
	if s.failure != nil {
		return nil, s.failure
	}

	saved := *v
	saved.ID = 12
	s.saved = &saved

	return &saved, nil
}

func (s *roleStore) Get(_ context.Context, _ int) (*domain.Role, error) { return s.saved, s.failure }

func (s *roleStore) List(_ context.Context, _ domain.ScopeRef) ([]*domain.Role, error) {
	return nil, s.failure
}

func (s *roleStore) Delete(_ context.Context, _ int, version int64) error {
	s.expected = version
	return s.failure
}

func TestHTTP(t *testing.T) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	store := &roleStore{}
	handler := api.NewHandler(application.NewRoles(store, catalog), nil, nil, application.NewPermissions(catalog))

	server, err := api.NewServer(handler, func(context.Context, *http.Request) (int, error) { return 37, nil })
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, method, path, body string
		status                   int
		failure                  error
	}{
		{"empty list", "GET", "/roles?scope_type=project&scope_id=1", "", 200, nil},
		{"bad id", "GET", "/roles/abc", "", 400, nil},
		{"negative id", "GET", "/roles/-1", "", 400, nil},
		{"bad scope", "GET", "/roles?scope_type=wrong&scope_id=1", "", 400, nil},
		{"missing filter", "GET", "/roles", "", 400, nil},
		{"bad json", "POST", "/roles", "{", 400, nil},
		{"missing body fields", "POST", "/roles", "{}", 400, nil},
		{"missing version", "DELETE", "/roles/1", "", 400, nil},
		{"delete version", "DELETE", "/roles/1?expected_version=7", "", 204, nil},
		{"not found", "GET", "/roles/1", "", 404, application.ErrNotFound},
		{"conflict", "DELETE", "/roles/1?expected_version=7", "", 409, application.ErrConflict},
		{"internal", "GET", "/roles/1", "", 500, errors.New("secret database detail")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store.failure = tc.failure
			request := httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")

			result := httptest.NewRecorder()
			server.ServeHTTP(result, request)
			validateResponse(t, request, result)

			if result.Code != tc.status {
				t.Fatalf("status %d: %s", result.Code, result.Body.String())
			}

			if strings.Contains(result.Body.String(), "secret") {
				t.Fatal("internal error exposed")
			}

			if tc.status >= 400 && !strings.Contains(result.Body.String(), `"code"`) {
				t.Fatal("missing error envelope")
			}

			if tc.name == "empty list" && !strings.Contains(result.Body.String(), `"items":[]`) {
				t.Fatal("empty list must be []")
			}
		})
	}

	if store.expected != 7 {
		t.Fatal("expected version was not forwarded")
	}

	store.failure = nil
	request := httptest.NewRequestWithContext(t.Context(), "POST", "/roles", strings.NewReader(`{"name":"reader","scope":{"type":"project","id":1},"permission_codes":["project.view"]}`))
	request.Header.Set("Content-Type", "application/json")

	result := httptest.NewRecorder()
	server.ServeHTTP(result, request)

	if result.Code != 201 {
		t.Fatalf("create: %d %s", result.Code, result.Body.String())
	}

	if store.saved.CreatedBy != 37 {
		t.Fatal("actor was not forwarded")
	}

	unavailable, err := api.NewServer(handler, nil)
	if err != nil {
		t.Fatal(err)
	}

	result = httptest.NewRecorder()
	unavailable.ServeHTTP(result, httptest.NewRequestWithContext(t.Context(), "GET", "/permissions?scope_type=project", nil))

	if result.Code != 503 {
		t.Fatal("unconfigured authentication must reject requests")
	}
}

func (s *roleStore) Update(_ context.Context, v *domain.Role, version int64) error {
	s.saved = v
	s.expected = version

	return s.failure
}

type templateStore struct {
	application.TemplateRepository
	saved    *domain.RoleTemplate
	expected int64
}

func (s *templateStore) Insert(_ context.Context, v *domain.RoleTemplate) (*domain.RoleTemplate, error) {
	saved := *v
	saved.ID = 8
	s.saved = &saved

	return &saved, nil
}

func (s *templateStore) Get(_ context.Context, _ int) (*domain.RoleTemplate, error) {
	return s.saved, nil
}

func (s *templateStore) List(_ context.Context, _ application.TemplateFilter) ([]*domain.RoleTemplate, error) {
	return []*domain.RoleTemplate{s.saved}, nil
}

func (s *templateStore) Update(_ context.Context, v *domain.RoleTemplate, version int64) error {
	s.saved = v
	s.expected = version

	return nil
}

func (s *templateStore) Delete(_ context.Context, _ int, version int64) error {
	s.expected = version
	return nil
}

type assignmentStore struct {
	application.AssignmentRepository
	saved          domain.RoleAssignment
	roleID, userID int
}

func (s *assignmentStore) Assign(_ context.Context, v domain.RoleAssignment) error {
	s.saved = v
	return nil
}

func (s *assignmentStore) Revoke(_ context.Context, roleID, userID int) error {
	s.roleID = roleID
	s.userID = userID

	return nil
}

func (s *assignmentStore) ListAssignments(_ context.Context, _ int) ([]domain.RoleAssignment, error) {
	return []domain.RoleAssignment{s.saved}, nil
}

func TestHTTPWorkflow(t *testing.T) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	r := &roleStore{}
	templates := &templateStore{}
	assignments := &assignmentStore{}

	server, err := api.NewServer(api.NewHandler(application.NewRoles(r, catalog), application.NewTemplates(templates, catalog), application.NewAssignments(assignments), application.NewPermissions(catalog)), func(context.Context, *http.Request) (int, error) { return 37, nil })
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"POST", "/role-templates", `{"name":"base","system_id":1,"target_scope_type":"project","permission_codes":["project.view"]}`, 201},
		{"GET", "/role-templates/8", "", 400},
		{"GET", "/role-templates/8?system_id=2", "", 404},
		{"PUT", "/role-templates/8?system_id=2", `{"name":"wrong","permission_codes":["project.view"],"expected_version":1}`, 404},
		{"DELETE", "/role-templates/8?system_id=2&expected_version=1", "", 404},
		{"GET", "/role-templates/8?system_id=1", "", 200},
		{"GET", "/role-templates?system_id=1&target_scope_type=project", "", 200},
		{"PUT", "/role-templates/8?system_id=1", `{"name":"changed","permission_codes":["project.view"],"expected_version":1}`, 200},
		{"POST", "/roles", `{"name":"role","scope":{"type":"project","id":1},"permission_codes":["project.view"]}`, 201},
		{"PUT", "/roles/12", `{"name":"changed","permission_codes":["project.view"],"expected_version":1}`, 200},
		{"GET", "/roles/12", "", 200},
		{"PUT", "/roles/12/assignments/9", "", 201},
		{"GET", "/roles/12/assignments", "", 200},
		{"GET", "/users/9/role-assignments?scope_type=project&scope_id=1", "", 200},
		{"GET", "/users/9/role-assignments?scope_id=1", "", 400},
		{"GET", "/users/9/role-assignments?scope_type=project", "", 400},
		{"GET", "/role-assignments?scope_type=project&scope_id=1", "", 200},
		{"GET", "/role-assignments", "", 400},
		{"DELETE", "/roles/12/assignments/9", "", 204},
		{"GET", "/permissions?scope_type=project", "", 200},
		{"DELETE", "/role-templates/8?system_id=1&expected_version=2", "", 204},
		{"DELETE", "/roles/12?expected_version=2", "", 204},
	} {
		req := httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")

		result := httptest.NewRecorder()
		server.ServeHTTP(result, req)
		validateResponse(t, req, result)

		if result.Code != tc.status {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, result.Code, result.Body.String())
		}
	}

	if templates.saved.CreatedBy != 37 || assignments.saved.AssignedBy() != 37 || assignments.saved.RoleID() != 12 || assignments.saved.UserID() != 9 || assignments.roleID != 12 || assignments.userID != 9 || templates.expected != 2 || r.expected != 2 {
		t.Fatal("actor, IDs or versions were not forwarded")
	}
}

func (s *roleStore) Page(_ context.Context, _ domain.ScopeRef, o application.ListOptions) (application.Page[*domain.Role], error) {
	return application.Page[*domain.Role]{Items: []*domain.Role{}, Limit: o.Limit, Offset: o.Offset}, s.failure
}

func (s *templateStore) Page(_ context.Context, _ application.TemplateFilter, o application.ListOptions) (application.Page[*domain.RoleTemplate], error) {
	return application.Page[*domain.RoleTemplate]{Items: []*domain.RoleTemplate{s.saved}, Total: 1, Limit: o.Limit, Offset: o.Offset}, nil
}

func (s *assignmentStore) Ensure(_ context.Context, a domain.RoleAssignment) (application.AssignmentResult, error) {
	s.saved = a
	return application.AssignmentResult{Assignment: a, Created: true}, nil
}

func (s *assignmentStore) Page(_ context.Context, f application.AssignmentFilter) (application.Page[application.AssignmentView], error) {
	scope, err := domain.NewScopeRef(domain.ScopeProject, 1)
	if err != nil {
		return application.Page[application.AssignmentView]{}, err
	}

	return application.Page[application.AssignmentView]{Items: []application.AssignmentView{{RoleID: s.saved.RoleID(), UserID: s.saved.UserID(), RoleName: "Role", Scope: scope, AssignedBy: s.saved.AssignedBy(), AssignedAt: s.saved.AssignedAt()}}, Total: 1, Limit: f.Limit, Offset: f.Offset}, nil
}
