package domain

import (
	"time"
)

func RestoreRoleAssignment(roleID, userID, assignedBy int, assignedAt time.Time) (RoleAssignment, error) {
	if roleID <= 0 || userID <= 0 || assignedBy <= 0 || assignedAt.IsZero() {
		return RoleAssignment{}, NewError(CodeInvalidInput, "invalid persisted assignment")
	}

	return RoleAssignment{roleID: roleID, userID: userID, assignedBy: assignedBy, assignedAt: assignedAt}, nil
}
