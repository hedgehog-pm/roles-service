package domain

type PermissionDescriptor struct {
	Code        PermissionCode
	ScopeType   ScopeType
	Title       string
	Description string
}
