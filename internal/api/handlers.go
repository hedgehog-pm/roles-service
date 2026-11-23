package api

import (
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	application "github.com/hedgehog-pm/roles-service/internal/service"
)

type Handler struct {
	RoleHandler
	RoleTemplateHandler
	RoleAssignmentHandler
	PermissionHandler
}

func NewHandler(roles *application.Roles, templates *application.Templates, assignments *application.Assignments, permissions *application.Permissions) *Handler {
	return &Handler{RoleHandler: RoleHandler{roles: roles, permissions: permissions}, RoleTemplateHandler: RoleTemplateHandler{templates: templates, permissions: permissions}, assignments: assignments, PermissionHandler: PermissionHandler{permissions: permissions}}
}

var _ generated.StrictServerInterface = (*Handler)(nil)
