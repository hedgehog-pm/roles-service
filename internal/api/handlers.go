package api

import (
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	"github.com/hedgehog-pm/roles-service/internal/usecase/createrole"
	"github.com/hedgehog-pm/roles-service/internal/usecase/createtemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroleassignments"
	"github.com/hedgehog-pm/roles-service/internal/usecase/listroles"
	"github.com/hedgehog-pm/roles-service/internal/usecase/permission"
	"github.com/hedgehog-pm/roles-service/internal/usecase/role"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roleassignment"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roletemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/updaterole"
	"github.com/hedgehog-pm/roles-service/internal/usecase/updatetemplate"
)

type Handler struct {
	RoleHandler
	RoleTemplateHandler
	RoleAssignmentHandler
	PermissionHandler
}

func NewHandler(roles RoleHandler, templates RoleTemplateHandler, assignments RoleAssignmentHandler, permissions *permission.Usecase) *Handler {
	return &Handler{RoleHandler: roles, RoleTemplateHandler: templates, RoleAssignmentHandler: assignments, PermissionHandler: PermissionHandler{permissions: permissions}}
}

func NewRoleHandler(roles *role.Usecase, list *listroles.Usecase, create *createrole.Usecase, update *updaterole.Usecase, permissions *permission.Usecase) RoleHandler {
	return RoleHandler{roles: roles, listroles: list, createrole: create, updaterole: update, permissions: permissions}
}

func NewRoleTemplateHandler(templates *roletemplate.Usecase, create *createtemplate.Usecase, update *updatetemplate.Usecase, permissions *permission.Usecase) RoleTemplateHandler {
	return RoleTemplateHandler{templates: templates, createtemplate: create, updatetemplate: update, permissions: permissions}
}

func NewRoleAssignmentHandler(assignments *roleassignment.Usecase, list *listroleassignments.Usecase) RoleAssignmentHandler {
	return RoleAssignmentHandler{assignments: assignments, listassignments: list}
}

var _ generated.StrictServerInterface = (*Handler)(nil)
