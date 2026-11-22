package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/hedgehog-pm/roles-service/internal/api"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api/roles"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

func validateResponse(t *testing.T, request *http.Request, result *httptest.ResponseRecorder) {
	t.Helper()

	spec, err := generated.GetSpec()
	if err != nil {
		t.Fatal(err)
	}

	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		t.Fatal(err)
	}

	route, params, err := router.FindRoute(request)
	if err != nil {
		t.Fatal(err)
	}

	input := &openapi3filter.ResponseValidationInput{RequestValidationInput: &openapi3filter.RequestValidationInput{Request: request, PathParams: params, Route: route}, Status: result.Code, Header: result.Header(), Options: &openapi3filter.Options{IncludeResponseStatus: true}}
	input.SetBodyBytes(result.Body.Bytes())

	if err := openapi3filter.ValidateResponse(t.Context(), input); err != nil {
		t.Fatal(err)
	}
}

func TestPermissionResponsesAndValidation(t *testing.T) {
	catalog, err := domain.NewPermissionCatalog()
	if err != nil {
		t.Fatal(err)
	}

	server, err := api.NewServer(api.NewHandler(application.NewRoles(&roleStore{}, catalog), nil, nil, application.NewPermissions(catalog)), func(context.Context, *http.Request) (int, error) { return 1, nil })
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		codes      []string
		name, code string
	}{
		{[]string{"unknown"}, "valid", "unknown_permission_code"},
		{[]string{string(domain.SystemRoleCreate)}, "valid", "permission_scope_mismatch"},
		{[]string{}, "valid", "empty_permissions"},
		{[]string{string(domain.ProjectView), string(domain.ProjectView)}, "valid", "invalid_input"},
		{[]string{string(domain.ProjectView)}, "   ", "invalid_input"},
		{[]string{string(domain.ProjectView)}, strings.Repeat("я", 201), "invalid_input"},
	} {
		data, err := json.Marshal(map[string]any{"name": tc.name, "scope": map[string]any{"type": "project", "id": 1}, "permission_codes": tc.codes})
		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequestWithContext(t.Context(), "POST", "/roles", strings.NewReader(string(data)))
		request.Header.Set("Content-Type", "application/json")

		result := httptest.NewRecorder()
		server.ServeHTTP(result, request)

		var body generated.Error
		if err := json.Unmarshal(result.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}

		if result.Code != 400 || body.Code != tc.code || body.Details == nil || len(*body.Details) == 0 {
			t.Fatalf("validation: %d %s", result.Code, result.Body.String())
		}

		validateResponse(t, request, result)
	}

	request := httptest.NewRequestWithContext(t.Context(), "POST", "/roles", strings.NewReader(`{"name":"Reader","scope":{"type":"project","id":1},"permission_codes":["project.view"]}`))
	request.Header.Set("Content-Type", "application/json")

	result := httptest.NewRecorder()
	server.ServeHTTP(result, request)
	validateResponse(t, request, result)

	var role generated.Role
	if err := json.Unmarshal(result.Body.Bytes(), &role); err != nil {
		t.Fatal(err)
	}

	if len(role.Permissions) != 1 || role.Permissions[0].Code != "project.view" || role.Permissions[0].Title == "" || role.Permissions[0].Description == "" || strings.Contains(result.Body.String(), "permission_codes") {
		t.Fatalf("permissions: %s", result.Body.String())
	}

	for _, scope := range []domain.ScopeType{domain.ScopeSystem, domain.ScopeWorkspace, domain.ScopeProject} {
		for _, entry := range catalog.List(scope) {
			if strings.TrimSpace(entry.Title) == "" || strings.TrimSpace(entry.Description) == "" {
				t.Errorf("missing catalog text: %s", entry.Code)
			}
		}
	}
}
