package api

import (
	"context"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
	application "github.com/hedgehog-pm/roles-service/internal/service"
)

type RoleTemplateHandler struct {
	templates   *application.Templates
	permissions *application.Permissions
}

func (h *RoleTemplateHandler) templateDTO(v *domain.RoleTemplate) generated.RoleTemplate {
	return generated.RoleTemplate{ID: v.ID, Name: v.Name, Permissions: permissionDTO(h.permissions, v.Permissions), Version: v.Version, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, SystemID: v.SystemID, TargetScopeType: generated.ScopeType(v.TargetScopeType)}
}

func (h *RoleTemplateHandler) CreateTemplate(ctx context.Context, request generated.CreateTemplateRequestObject) (generated.CreateTemplateResponseObject, error) {
	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	value, err := h.templates.Create(ctx, actor, application.CreateTemplateInput{Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes), SystemID: request.Body.SystemID, TargetScopeType: domain.ScopeType(request.Body.TargetScopeType)})
	if err != nil {
		return nil, err
	}

	return generated.CreateTemplate201JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) GetTemplate(ctx context.Context, request generated.GetTemplateRequestObject) (generated.GetTemplateResponseObject, error) {
	if err := h.checkTemplateSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	value, err := h.templates.Get(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return generated.GetTemplate200JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) UpdateTemplate(ctx context.Context, request generated.UpdateTemplateRequestObject) (generated.UpdateTemplateResponseObject, error) {
	if err := h.checkTemplateSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	value, err := h.templates.Update(ctx, actor, application.UpdateTemplateInput{ID: request.ID, Version: request.Body.ExpectedVersion, Name: request.Body.Name, PermissionCodes: permissionCodes(request.Body.PermissionCodes)})
	if err != nil {
		return nil, err
	}

	return generated.UpdateTemplate200JSONResponse(h.templateDTO(value)), nil
}

func (h *RoleTemplateHandler) DeleteTemplate(ctx context.Context, request generated.DeleteTemplateRequestObject) (generated.DeleteTemplateResponseObject, error) {
	if err := h.checkTemplateSystem(ctx, request.ID, request.Params.SystemID); err != nil {
		return nil, err
	}

	actor, err := actorID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.templates.Delete(ctx, actor, request.ID, request.Params.ExpectedVersion); err != nil {
		return nil, err
	}

	return generated.DeleteTemplate204Response{}, nil
}

func (h *RoleTemplateHandler) ListTemplates(ctx context.Context, request generated.ListTemplatesRequestObject) (generated.ListTemplatesResponseObject, error) {
	filter := domain.TemplateFilter{SystemID: request.Params.SystemID, TargetScopeType: domain.ScopeType(request.Params.TargetScopeType)}

	options := domain.ListOptions{Pagination: pagination(request.Params.Limit, request.Params.Offset), Search: optional(request.Params.Search), Sort: string(optional(request.Params.Sort)), Direction: string(optional(request.Params.Direction))}

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

func (h *RoleTemplateHandler) checkTemplateSystem(ctx context.Context, id, systemID int) error {
	value, err := h.templates.Get(ctx, id)
	if err != nil {
		return err
	}

	if value.SystemID != systemID {
		return domain.ErrNotFound
	}

	return nil
}
