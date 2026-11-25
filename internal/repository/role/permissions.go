package role

import "github.com/hedgehog-pm/roles-service/internal/domain"

func permissionCodes(values []string) []domain.PermissionCode {
	codes := make([]domain.PermissionCode, len(values))
	for i, value := range values {
		codes[i] = domain.PermissionCode(value)
	}

	return codes
}
