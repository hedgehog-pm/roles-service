package api

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroleassignments"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roleassignment"
)

type RoleAssignmentHandler struct {
	listassignments *listroleassignments.Usecase
	assignments     *roleassignment.Usecase
}

func (h *RoleAssignmentHandler) AssignRole(ctx context.Context, request generated.AssignRoleRequestObject) (generated.AssignRoleResponseObject, error) {
	actor := actorID(ctx)

	result, err := h.assignments.Assign(ctx, actor, request.RoleID, request.UserID)
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
	actor := actorID(ctx)

	if err := h.assignments.Revoke(ctx, actor, request.RoleID, request.UserID); err != nil {
		return nil, err
	}

	return generated.RevokeRole204Response{}, nil
}

func (h *RoleAssignmentHandler) ListAssignments(ctx context.Context, request generated.ListAssignmentsRequestObject) (generated.ListAssignmentsResponseObject, error) {
	values, err := h.listassignments.Do(ctx, listroleassignments.Params{RoleID: request.RoleID, Pagination: pageParams(request.Params.Limit, request.Params.Offset)})
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
	values, err := h.listassignments.Do(ctx, listroleassignments.Params{UserID: request.UserID, ScopeType: domain.ScopeType(request.Params.ScopeType), ScopeID: request.Params.ScopeID, Pagination: pageParams(request.Params.Limit, request.Params.Offset)})
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
	values, err := h.listassignments.Do(ctx, listroleassignments.Params{ScopeType: domain.ScopeType(request.Params.ScopeType), ScopeID: request.Params.ScopeID, Pagination: pageParams(request.Params.Limit, request.Params.Offset)})
	if err != nil {
		return nil, err
	}

	result := generated.ListScopeAssignments200JSONResponse{Items: make([]generated.UserRoleAssignment, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, generated.UserRoleAssignment{RoleID: v.RoleID, UserID: v.UserID, RoleName: v.RoleName, Scope: generated.Scope{ID: v.Scope.ID(), Type: generated.ScopeType(v.Scope.Type())}, AssignedBy: v.AssignedBy, AssignedAt: v.AssignedAt})
	}

	return result, nil
}
