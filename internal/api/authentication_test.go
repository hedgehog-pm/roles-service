package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
)

func TestAuthentication(t *testing.T) {
	tests := []struct {
		name         string
		headers      []string
		unconfigured bool
		id           int
		err          error
		status       int
		code         string
		verifyCalled bool
	}{
		{name: "not configured", unconfigured: true, status: http.StatusServiceUnavailable, code: "unavailable"},
		{name: "missing header", status: http.StatusUnauthorized, code: "unauthorized"},
		{name: "empty token", headers: []string{"Bearer "}, status: http.StatusUnauthorized, code: "unauthorized"},
		{name: "wrong scheme", headers: []string{"Basic token"}, status: http.StatusUnauthorized, code: "unauthorized"},
		{name: "extra fields", headers: []string{"Bearer token extra"}, status: http.StatusUnauthorized, code: "unauthorized"},
		{name: "duplicate headers", headers: []string{"Bearer token", "Bearer other"}, status: http.StatusUnauthorized, code: "unauthorized"},
		{name: "rejected token", headers: []string{"Bearer token"}, err: fmt.Errorf("expired: %w", ErrUnauthenticated), status: http.StatusUnauthorized, code: "unauthorized", verifyCalled: true},
		{name: "service failure", headers: []string{"Bearer token"}, err: errors.New("private upstream details"), status: http.StatusServiceUnavailable, code: "unavailable", verifyCalled: true},
		{name: "zero user", headers: []string{"Bearer token"}, status: http.StatusUnauthorized, code: "unauthorized", verifyCalled: true},
		{name: "negative user", headers: []string{"Bearer token"}, id: -1, status: http.StatusUnauthorized, code: "unauthorized", verifyCalled: true},
		{name: "valid token", headers: []string{"Bearer token"}, id: 42, status: http.StatusOK, verifyCalled: true},
		{name: "case insensitive scheme", headers: []string{"bearer token"}, id: 42, status: http.StatusOK, verifyCalled: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &permissionHandlerStub{}

			server, err := NewServer(handler)
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/permissions?scope_type=project", nil)
			for _, header := range tt.headers {
				request.Header.Add("Authorization", header)
			}

			called := false

			var verify VerifyToken
			if !tt.unconfigured {
				verify = func(ctx context.Context, token string) (int, error) {
					called = true

					if ctx != request.Context() || token != "token" {
						t.Fatal("verifier did not receive request context and raw token")
					}

					return tt.id, tt.err
				}
			}

			server.Use(Authentication(verify))

			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			if response.Code != tt.status || called != tt.verifyCalled {
				t.Fatalf("status=%d verifier called=%v body=%s", response.Code, called, response.Body.String())
			}

			if tt.status == http.StatusOK {
				if !handler.called || handler.actor != tt.id {
					t.Fatalf("handler called=%v actor=%d", handler.called, handler.actor)
				}
			} else {
				if handler.called {
					t.Fatal("rejected request reached handler")
				}

				var body generated.Error
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}

				if body.Code != tt.code || (body.Message != "unauthorized" && body.Message != "service unavailable") {
					t.Fatalf("unexpected public error: %+v", body)
				}
			}
		})
	}
}
