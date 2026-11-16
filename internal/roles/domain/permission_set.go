package domain

import "slices"

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
	if s.scopeType != other.scopeType {
		return false
	}

	if len(s.codes) != len(other.codes) {
		return false
	}

	for code := range s.codes {
		if !other.Contains(code) {
			return false
		}
	}

	return true
}
