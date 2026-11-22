package api

import (
	"context"

	generated "github.com/hedgehog-pm/roles-service/internal/generated/api/roles"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type PermissionHandler struct{ permissions *application.Permissions }

func (h *PermissionHandler) ListPermissions(_ context.Context, request generated.ListPermissionsRequestObject) (generated.ListPermissionsResponseObject, error) {
	values, err := h.permissions.List(domain.ScopeType(request.Params.ScopeType))
	if err != nil {
		return nil, err
	}

	result := make(generated.ListPermissions200JSONResponse, 0, len(values))
	for _, v := range values {
		result = append(result, generated.Permission{Code: string(v.Code), ScopeType: generated.ScopeType(v.ScopeType), Title: v.Title, Description: v.Description})
	}

	return result, nil
}
