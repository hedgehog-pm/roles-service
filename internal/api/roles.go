package api

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	application "github.com/hedgehog-pm/roles-service/internal/service"
)

type RoleHandler struct {
	roles       *application.Roles
	permissions *application.Permissions
}

func (h *RoleHandler) roleDTO(v *domain.Role) generated.Role {
	return generated.Role{ID: v.ID, Name: v.Name, Permissions: permissionDTO(h.permissions, v.Permissions), Version: v.Version, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, Scope: generated.Scope{ID: v.Scope.ID(), Type: generated.ScopeType(v.Scope.Type())}}
}

func (h *RoleHandler) CreateRole(ctx context.Context, request generated.CreateRoleRequestObject) (generated.CreateRoleResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	scope, err := domain.NewScopeRef(domain.ScopeType(request.Body.Scope.Type), request.Body.Scope.ID)
	if err != nil {
		return nil, err
	}

	value, err := h.roles.Create(ctx, actor, application.CreateRoleInput{Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes), Scope: scope})
	if err != nil {
		return nil, err
	}

	return generated.CreateRole201JSONResponse(h.roleDTO(value)), nil
}

func (h *RoleHandler) GetRole(ctx context.Context, request generated.GetRoleRequestObject) (generated.GetRoleResponseObject, error) {
	value, err := h.roles.Get(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return generated.GetRole200JSONResponse(h.roleDTO(value)), nil
}

func (h *RoleHandler) UpdateRole(ctx context.Context, request generated.UpdateRoleRequestObject) (generated.UpdateRoleResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	value, err := h.roles.Update(ctx, actor, application.UpdateRoleInput{ID: request.ID, Version: request.Body.ExpectedVersion, Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes)})
	if err != nil {
		return nil, err
	}

	return generated.UpdateRole200JSONResponse(h.roleDTO(value)), nil
}

func (h *RoleHandler) DeleteRole(ctx context.Context, request generated.DeleteRoleRequestObject) (generated.DeleteRoleResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.roles.Delete(ctx, actor, request.ID, request.Params.ExpectedVersion); err != nil {
		return nil, err
	}

	return generated.DeleteRole204Response{}, nil
}

func (h *RoleHandler) ListRoles(ctx context.Context, request generated.ListRolesRequestObject) (generated.ListRolesResponseObject, error) {
	filter, err := domain.NewScopeRef(domain.ScopeType(request.Params.ScopeType), request.Params.ScopeID)
	if err != nil {
		return nil, err
	}

	options := domain.ListOptions{Pagination: pagination(request.Params.Limit, request.Params.Offset), Search: optional(request.Params.Search), Sort: string(optional(request.Params.Sort)), Direction: string(optional(request.Params.Direction)), PermissionCode: domain.PermissionCode(optional(request.Params.PermissionCode))}

	values, err := h.roles.Page(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	result := generated.ListRoles200JSONResponse{Items: make([]generated.Role, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, h.roleDTO(v))
	}

	return result, nil
}
