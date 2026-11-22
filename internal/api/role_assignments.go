package api

import (
	"context"

	generated "github.com/hedgehog-pm/roles-service/internal/generated/api/roles"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

type RoleAssignmentHandler struct{ assignments *application.Assignments }

func (h *RoleAssignmentHandler) AssignRole(ctx context.Context, request generated.AssignRoleRequestObject) (generated.AssignRoleResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.assignments.Ensure(ctx, actor, request.RoleID, request.UserID)
	if err != nil {
		return nil, err
	}

	value := assignmentDTO(result.Assignment)
	if result.Created {
		return generated.AssignRole201JSONResponse(value), nil
	}

	return generated.AssignRole200JSONResponse(value), nil
}

func (h *RoleAssignmentHandler) RevokeRole(ctx context.Context, request generated.RevokeRoleRequestObject) (generated.RevokeRoleResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.assignments.Revoke(ctx, actor, request.RoleID, request.UserID); err != nil {
		return nil, err
	}

	return generated.RevokeRole204Response{}, nil
}

func (h *RoleAssignmentHandler) ListAssignments(ctx context.Context, request generated.ListAssignmentsRequestObject) (generated.ListAssignmentsResponseObject, error) {
	values, err := h.assignments.Page(ctx, application.AssignmentFilter{RoleID: request.RoleID, Pagination: pagination(request.Params.Limit, request.Params.Offset)})
	if err != nil {
		return nil, err
	}

	result := generated.ListAssignments200JSONResponse{Items: make([]generated.RoleAssignment, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, generated.RoleAssignment{RoleID: v.RoleID, UserID: v.UserID, AssignedBy: v.AssignedBy, AssignedAt: v.AssignedAt})
	}

	return result, nil
}

func (h *RoleAssignmentHandler) ListUserAssignments(ctx context.Context, request generated.ListUserAssignmentsRequestObject) (generated.ListUserAssignmentsResponseObject, error) {
	values, err := h.assignments.Page(ctx, application.AssignmentFilter{UserID: request.UserID, ScopeType: domain.ScopeType(request.Params.ScopeType), ScopeID: request.Params.ScopeID, Pagination: pagination(request.Params.Limit, request.Params.Offset)})
	if err != nil {
		return nil, err
	}

	result := generated.ListUserAssignments200JSONResponse{Items: make([]generated.UserRoleAssignment, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, generated.UserRoleAssignment{RoleID: v.RoleID, UserID: v.UserID, RoleName: v.RoleName, Scope: generated.Scope{ID: v.Scope.ID(), Type: generated.ScopeType(v.Scope.Type())}, AssignedBy: v.AssignedBy, AssignedAt: v.AssignedAt})
	}

	return result, nil
}

func (h *RoleAssignmentHandler) ListScopeAssignments(ctx context.Context, request generated.ListScopeAssignmentsRequestObject) (generated.ListScopeAssignmentsResponseObject, error) {
	values, err := h.assignments.Page(ctx, application.AssignmentFilter{ScopeType: domain.ScopeType(request.Params.ScopeType), ScopeID: request.Params.ScopeID, Pagination: pagination(request.Params.Limit, request.Params.Offset)})
	if err != nil {
		return nil, err
	}

	result := generated.ListScopeAssignments200JSONResponse{Items: make([]generated.UserRoleAssignment, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, generated.UserRoleAssignment{RoleID: v.RoleID, UserID: v.UserID, RoleName: v.RoleName, Scope: generated.Scope{ID: v.Scope.ID(), Type: generated.ScopeType(v.Scope.Type())}, AssignedBy: v.AssignedBy, AssignedAt: v.AssignedAt})
	}

	return result, nil
}
