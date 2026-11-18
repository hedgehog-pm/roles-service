package domain

import (
	"errors"

	"github.com/google/uuid"
)

type ScopeType string

const (
	ScopeSystem    ScopeType = "system"
	ScopeWorkspace ScopeType = "workspace"
	ScopeProject   ScopeType = "project"
)

// value obj
type ScopeRef struct {
	scopeType  ScopeType
	resourceID uuid.UUID
}

func NewScopeRef(scopeType ScopeType, id uuid.UUID) (ScopeRef, error) {
	switch scopeType {
	case ScopeSystem, ScopeWorkspace, ScopeProject:
	default:
		return ScopeRef{}, errors.New("unknown scope type")
	}

	if id == uuid.Nil {
		return ScopeRef{}, errors.New("scope ID is empty")
	}

	return ScopeRef{scopeType: scopeType, resourceID: id}, nil
}

func (s ScopeRef) Type() ScopeType {
	return s.scopeType
}

func (s ScopeRef) ID() uuid.UUID {
	return s.resourceID
}
