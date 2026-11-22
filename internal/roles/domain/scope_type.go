package domain

type ScopeType string

const (
	ScopeSystem    ScopeType = "system"
	ScopeWorkspace ScopeType = "workspace"
	ScopeProject   ScopeType = "project"
)

// value obj
type ScopeRef struct {
	scopeType  ScopeType
	resourceID int
}

func NewScopeRef(scopeType ScopeType, id int) (ScopeRef, error) {
	switch scopeType {
	case ScopeSystem, ScopeWorkspace, ScopeProject:
	default:
		return ScopeRef{}, NewError(CodeInvalidInput, "unknown scope type")
	}

	if id <= 0 {
		return ScopeRef{}, NewError(CodeInvalidInput, "scope ID must be positive")
	}

	return ScopeRef{scopeType: scopeType, resourceID: id}, nil
}

func (s ScopeRef) Type() ScopeType {
	return s.scopeType
}

func (s ScopeRef) ID() int {
	return s.resourceID
}
