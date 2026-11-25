package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
)

type permissionHandlerStub struct {
	generated.StrictServerInterface
	actor  int
	called bool
}

func (h *permissionHandlerStub) ListPermissions(ctx context.Context, _ generated.ListPermissionsRequestObject) (generated.ListPermissionsResponseObject, error) {
	h.called = true
	h.actor = actorID(ctx)

	return generated.ListPermissions200JSONResponse{}, nil
}

func TestServerDevelopmentActor(t *testing.T) {
	for _, authorization := range []string{"", "Bearer invalid-token"} {
		t.Run(authorization, func(t *testing.T) {
			handler := &permissionHandlerStub{}

			server, err := NewServer(handler)
			if err != nil {
				t.Fatal(err)
			}

			server.Use(DevelopmentActor())

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/permissions?scope_type=project", nil)
			request.Header.Set("Authorization", authorization)

			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			if response.Code != http.StatusOK || !handler.called || handler.actor != 1 {
				t.Fatalf("status=%d called=%v actor=%d body=%s", response.Code, handler.called, handler.actor, response.Body.String())
			}
		})
	}
}

func TestServerRequestValidation(t *testing.T) {
	for _, target := range []string{"/permissions", "/permissions?scope_type=unknown"} {
		t.Run(target, func(t *testing.T) {
			handler := &permissionHandlerStub{}

			server, err := NewServer(handler)
			if err != nil {
				t.Fatal(err)
			}

			server.Use(DevelopmentActor())

			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, nil))

			if response.Code != http.StatusBadRequest || handler.called || !strings.Contains(response.Body.String(), `"code":"invalid_input"`) {
				t.Fatalf("status=%d called=%v body=%s", response.Code, handler.called, response.Body.String())
			}
		})
	}
}
