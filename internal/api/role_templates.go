package api

import (
	"context"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
	"github.com/hedgehog-pm/roles-service/internal/usecase/createtemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/permission"
	"github.com/hedgehog-pm/roles-service/internal/usecase/roletemplate"
	"github.com/hedgehog-pm/roles-service/internal/usecase/updatetemplate"
)

type RoleTemplateHandler struct {
	createtemplate *createtemplate.Usecase
	updatetemplate *updatetemplate.Usecase
	templates      *roletemplate.Usecase
	permissions    *permission.Usecase
}

func (h *RoleTemplateHandler) templateDTO(v *domain.RoleTemplate) generated.RoleTemplate {
	return generated.RoleTemplate{ID: v.ID, Name: v.Name, Permissions: permissionDTO(h.permissions, v.Permissions), Version: v.Version, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, SystemID: v.SystemID, TargetScopeType: generated.ScopeType(v.TargetScopeType)}
}

func (h *RoleTemplateHandler) CreateTemplate(ctx context.Context, request generated.CreateTemplateRequestObject) (generated.CreateTemplateResponseObject, error) {
	actor := actorID(ctx)

	value, err := h.createtemplate.Do(ctx, actor, createtemplate.Params{Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes), SystemID: request.Body.SystemID, TargetScopeType: domain.ScopeType(request.Body.TargetScopeType)})
	if err != nil {
		return nil, err
	}

	return generated.CreateTemplate201JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) GetTemplate(ctx context.Context, request generated.GetTemplateRequestObject) (generated.GetTemplateResponseObject, error) {
	if err := h.templates.CheckSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	value, err := h.templates.Get(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return generated.GetTemplate200JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) UpdateTemplate(ctx context.Context, request generated.UpdateTemplateRequestObject) (generated.UpdateTemplateResponseObject, error) {
	if err := h.templates.CheckSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	actor := actorID(ctx)

	value, err := h.updatetemplate.Do(ctx, actor, updatetemplate.Params{ID: request.ID, Version: request.Body.ExpectedVersion, Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes)})
	if err != nil {
		return nil, err
	}

	return generated.UpdateTemplate200JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) DeleteTemplate(ctx context.Context, request generated.DeleteTemplateRequestObject) (generated.DeleteTemplateResponseObject, error) {
	if err := h.templates.CheckSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	actor := actorID(ctx)

	if err := h.templates.Delete(ctx, actor, request.ID, request.Params.ExpectedVersion); err != nil {
		return nil, err
	}

	return generated.DeleteTemplate204Response{}, nil
}

func (h *RoleTemplateHandler) ListTemplates(ctx context.Context, request generated.ListTemplatesRequestObject) (generated.ListTemplatesResponseObject, error) {
	filter := roletemplate.Filter{SystemID: request.Params.SystemID, TargetScopeType: domain.ScopeType(request.Params.TargetScopeType)}

	options := pagination.Options{Pagination: pageParams(request.Params.Limit, request.Params.Offset), Search: lo.FromPtr(request.Params.Search), Sort: string(lo.FromPtr(request.Params.Sort)), Direction: string(lo.FromPtr(request.Params.Direction))}

	values, err := h.templates.Page(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	result := generated.ListTemplates200JSONResponse{Items: make([]generated.RoleTemplate, 0, len(values.Items)), Total: values.Total, Limit: values.Limit, Offset: values.Offset}
	for _, v := range values.Items {
		result.Items = append(result.Items, h.templateDTO(v))
	}

	return result, nil
}
