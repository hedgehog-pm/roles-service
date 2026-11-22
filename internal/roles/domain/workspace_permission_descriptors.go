package domain

func workspacePermissionDescriptors() []PermissionDescriptor {
	return []PermissionDescriptor{
		{
			Code:        WorkspaceProjectCreate,
			ScopeType:   ScopeWorkspace,
			Title:       "Создание проектов",
			Description: "Разрешает создавать проекты в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceProjectView,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр проектов",
			Description: "Разрешает просматривать список доступных проектов рабочего пространства.",
		},
		{
			Code:        WorkspaceProjectViewAll,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр всех проектов",
			Description: "Разрешает просматривать все проекты рабочего пространства, включая скрытые и проекты без участия пользователя.",
		},
		{
			Code:        WorkspaceProjectEditAll,
			ScopeType:   ScopeWorkspace,
			Title:       "Редактирование всех проектов",
			Description: "Разрешает редактировать все проекты рабочего пространства, включая проекты без участия пользователя; доступ к сущностям наследуется согласно политике проекта.",
		},
		{
			Code:        WorkspaceProjectDeleteAll,
			ScopeType:   ScopeWorkspace,
			Title:       "Удаление всех проектов",
			Description: "Разрешает удалять любые проекты рабочего пространства, включая проекты без участия пользователя.",
		},
		{
			Code:        WorkspaceRoleView,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр ролей рабочего пространства",
			Description: "Разрешает просматривать роли рабочего пространства в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceRoleCreate,
			ScopeType:   ScopeWorkspace,
			Title:       "Создание ролей рабочего пространства",
			Description: "Разрешает создавать роли рабочего пространства в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceRoleEdit,
			ScopeType:   ScopeWorkspace,
			Title:       "Редактирование ролей рабочего пространства",
			Description: "Разрешает редактировать роли рабочего пространства в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceRoleDelete,
			ScopeType:   ScopeWorkspace,
			Title:       "Удаление ролей рабочего пространства",
			Description: "Разрешает удалять роли рабочего пространства в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceRoleAssign,
			ScopeType:   ScopeWorkspace,
			Title:       "Назначение ролей рабочего пространства",
			Description: "Разрешает назначать роли рабочего пространства пользователям в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceRoleUnassign,
			ScopeType:   ScopeWorkspace,
			Title:       "Снятие ролей рабочего пространства",
			Description: "Разрешает снимать роли рабочего пространства с пользователей в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceParticipantView,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр участников рабочего пространства",
			Description: "Разрешает просматривать участников рабочего пространства в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceParticipantAdd,
			ScopeType:   ScopeWorkspace,
			Title:       "Добавление участников рабочего пространства",
			Description: "Разрешает добавлять пользователей в участники рабочего пространства.",
		},
		{
			Code:        WorkspaceParticipantRemove,
			ScopeType:   ScopeWorkspace,
			Title:       "Удаление участников рабочего пространства",
			Description: "Разрешает удалять пользователей из участников рабочего пространства.",
		},
		{
			Code:        WorkspaceInvitationCreate,
			ScopeType:   ScopeWorkspace,
			Title:       "Создание приглашений в рабочее пространство",
			Description: "Разрешает создавать приглашения в рабочее пространство в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceInvitationView,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр приглашений в рабочее пространство",
			Description: "Разрешает просматривать приглашения в рабочее пространство в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceInvitationEdit,
			ScopeType:   ScopeWorkspace,
			Title:       "Редактирование приглашений в рабочее пространство",
			Description: "Разрешает редактировать приглашения в рабочее пространство в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceInvitationDelete,
			ScopeType:   ScopeWorkspace,
			Title:       "Удаление приглашений в рабочее пространство",
			Description: "Разрешает удалять приглашения в рабочее пространство в текущем рабочем пространстве.",
		},
		{
			Code:        WorkspaceView,
			ScopeType:   ScopeWorkspace,
			Title:       "Просмотр рабочего пространства",
			Description: "Разрешает просматривать рабочее пространство.",
		},
		{
			Code:        WorkspaceEdit,
			ScopeType:   ScopeWorkspace,
			Title:       "Редактирование рабочего пространства",
			Description: "Разрешает редактировать рабочее пространство.",
		},
		{
			Code:        WorkspaceDelete,
			ScopeType:   ScopeWorkspace,
			Title:       "Удаление рабочего пространства",
			Description: "Разрешает удалять рабочее пространство.",
		},
	}
}
