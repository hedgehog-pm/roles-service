package domain

func systemPermissionDescriptors() []PermissionDescriptor {
	return []PermissionDescriptor{
		{
			Code:        SystemWorkspaceCreate,
			ScopeType:   ScopeSystem,
			Title:       "Создание рабочих пространств",
			Description: "Разрешает создавать рабочие пространства в системе.",
		},
		{
			Code:        SystemRoleCreate,
			ScopeType:   ScopeSystem,
			Title:       "Создание системных ролей",
			Description: "Разрешает создавать системные роли в системе.",
		},
		{
			Code:        SystemRoleEdit,
			ScopeType:   ScopeSystem,
			Title:       "Редактирование системных ролей",
			Description: "Разрешает редактировать системные роли в системе.",
		},
		{
			Code:        SystemRoleView,
			ScopeType:   ScopeSystem,
			Title:       "Просмотр системных ролей",
			Description: "Разрешает просматривать системные роли в системе.",
		},
		{
			Code:        SystemRoleDelete,
			ScopeType:   ScopeSystem,
			Title:       "Удаление системных ролей",
			Description: "Разрешает удалять системные роли в системе.",
		},
		{
			Code:        SystemRoleAssign,
			ScopeType:   ScopeSystem,
			Title:       "Назначение системных ролей",
			Description: "Разрешает назначать системные роли пользователям в системе.",
		},
		{
			Code:        SystemRoleUnassign,
			ScopeType:   ScopeSystem,
			Title:       "Снятие системных ролей",
			Description: "Разрешает снимать системные роли с пользователей в системе.",
		},
		{
			Code:        SystemWorkspaceRoleTemplateCreate,
			ScopeType:   ScopeSystem,
			Title:       "Создание шаблонов ролей рабочих пространств",
			Description: "Разрешает создавать шаблоны ролей рабочих пространств в системе.",
		},
		{
			Code:        SystemWorkspaceRoleTemplateView,
			ScopeType:   ScopeSystem,
			Title:       "Просмотр шаблонов ролей рабочих пространств",
			Description: "Разрешает просматривать шаблоны ролей рабочих пространств в системе.",
		},
		{
			Code:        SystemWorkspaceRoleTemplateEdit,
			ScopeType:   ScopeSystem,
			Title:       "Редактирование шаблонов ролей рабочих пространств",
			Description: "Разрешает редактировать шаблоны ролей рабочих пространств в системе.",
		},
		{
			Code:        SystemWorkspaceRoleTemplateDelete,
			ScopeType:   ScopeSystem,
			Title:       "Удаление шаблонов ролей рабочих пространств",
			Description: "Разрешает удалять шаблоны ролей рабочих пространств в системе.",
		},
		{
			Code:        SystemProjectRoleTemplateCreate,
			ScopeType:   ScopeSystem,
			Title:       "Создание шаблонов проектных ролей",
			Description: "Разрешает создавать шаблоны проектных ролей в системе.",
		},
		{
			Code:        SystemProjectRoleTemplateView,
			ScopeType:   ScopeSystem,
			Title:       "Просмотр шаблонов проектных ролей",
			Description: "Разрешает просматривать шаблоны проектных ролей в системе.",
		},
		{
			Code:        SystemProjectRoleTemplateEdit,
			ScopeType:   ScopeSystem,
			Title:       "Редактирование шаблонов проектных ролей",
			Description: "Разрешает редактировать шаблоны проектных ролей в системе.",
		},
		{
			Code:        SystemProjectRoleTemplateDelete,
			ScopeType:   ScopeSystem,
			Title:       "Удаление шаблонов проектных ролей",
			Description: "Разрешает удалять шаблоны проектных ролей в системе.",
		},
	}
}
