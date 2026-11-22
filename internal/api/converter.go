package api

import (
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api/roles"
	"github.com/hedgehog-pm/roles-service/internal/roles/application"
	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

func permissionCodes(values []string) []domain.PermissionCode {
	result := make([]domain.PermissionCode, len(values))
	for i, v := range values {
		result[i] = domain.PermissionCode(v)
	}

	return result
}

func permissionDTO(permissions *application.Permissions, set domain.PermissionSet) []generated.Permission {
	values := permissions.Describe(set)

	result := make([]generated.Permission, 0, len(values))
	for _, v := range values {
		result = append(result, generated.Permission{Code: string(v.Code), ScopeType: generated.ScopeType(v.ScopeType), Title: v.Title, Description: v.Description})
	}

	return result
}

func assignmentDTO(v domain.RoleAssignment) generated.RoleAssignment {
	return generated.RoleAssignment{RoleID: v.RoleID(), UserID: v.UserID(), AssignedBy: v.AssignedBy(), AssignedAt: v.AssignedAt()}
}
