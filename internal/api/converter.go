package api

import (
	"github.com/samber/lo"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	"github.com/hedgehog-pm/roles-service/internal/usecase/permission"
)

func permissionCodes(values []string) []domain.PermissionCode {
	return lo.Map(values, func(v string, _ int) domain.PermissionCode { return domain.PermissionCode(v) })
}

func permissionDTO(permissions *permission.Usecase, set domain.PermissionSet) []generated.Permission {
	values := permissions.Describe(set)

	return lo.Map(values, func(v domain.PermissionDescriptor, _ int) generated.Permission {
		return generated.Permission{Code: string(v.Code), ScopeType: generated.ScopeType(v.ScopeType), Title: v.Title, Description: v.Description}
	})
}

func assignmentDTO(v domain.RoleAssignment) generated.RoleAssignment {
	return generated.RoleAssignment{RoleID: v.RoleID(), UserID: v.UserID(), AssignedBy: v.AssignedBy(), AssignedAt: v.AssignedAt()}
}
