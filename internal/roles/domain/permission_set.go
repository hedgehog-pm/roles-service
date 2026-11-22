package domain

import (
	"maps"
	"slices"
)

type PermissionSet struct {
	scopeType ScopeType
	codes     map[PermissionCode]struct{}
}

func (s PermissionSet) ScopeType() ScopeType {
	return s.scopeType
}

func (s PermissionSet) Contains(code PermissionCode) bool {
	_, exists := s.codes[code]
	return exists
}

func (s PermissionSet) Codes() []PermissionCode {
	result := make([]PermissionCode, 0, len(s.codes))

	for code := range s.codes {
		result = append(result, code)
	}

	slices.Sort(result)

	return result
}

func (s PermissionSet) Equal(other PermissionSet) bool {
	return s.scopeType == other.scopeType && maps.Equal(s.codes, other.codes)
}
