package domain

import (
	"fmt"
	"sort"
)

type PermissionCatalog struct {
	entries map[PermissionCode]PermissionDescriptor
}

func NewPermissionCatalog() (PermissionCatalog, error) {
	groups := [][]PermissionDescriptor{
		systemPermissionDescriptors(),
		workspacePermissionDescriptors(),
		projectPermissionDescriptors(),
	}

	entries := make(map[PermissionCode]PermissionDescriptor)

	for _, descriptors := range groups {
		for _, descriptor := range descriptors {
			if _, exists := entries[descriptor.Code]; exists {
				return PermissionCatalog{}, fmt.Errorf(
					"duplicate permission code: %s",
					descriptor.Code,
				)
			}

			entries[descriptor.Code] = descriptor
		}
	}

	return PermissionCatalog{entries: entries}, nil
}

func (c PermissionCatalog) LookUp(code PermissionCode) (PermissionDescriptor, bool) {
	descriptor, exists := c.entries[code]
	return descriptor, exists
}

func (c PermissionCatalog) List(scopeType ScopeType) []PermissionDescriptor {
	result := make([]PermissionDescriptor, 0)

	for _, descriptor := range c.entries {
		if descriptor.ScopeType == scopeType {
			result = append(result, descriptor)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result
}

func (c PermissionCatalog) NewSet(
	scopeType ScopeType,
	codes []PermissionCode,
) (PermissionSet, error) {
	switch scopeType {
	case ScopeSystem, ScopeWorkspace, ScopeProject:
	default:
		return PermissionSet{}, fmt.Errorf(
			"unknown scope type: %s",
			scopeType,
		)
	}

	uniqueCodes := make(map[PermissionCode]struct{}, len(codes))

	for _, code := range codes {
		descriptor, exists := c.entries[code]
		if !exists {
			return PermissionSet{}, fmt.Errorf(
				"unknown permission code: %s",
				code,
			)
		}

		if descriptor.ScopeType != scopeType {
			return PermissionSet{}, fmt.Errorf(
				"permission %s belongs to scope %s, not %s",
				code,
				descriptor.ScopeType,
				scopeType,
			)
		}

		uniqueCodes[code] = struct{}{}
	}

	return PermissionSet{
		scopeType: scopeType,
		codes:     uniqueCodes,
	}, nil
}
