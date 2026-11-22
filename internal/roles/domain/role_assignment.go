package domain

import (
	"time"
)

type RoleAssignment struct {
	roleID     int
	userID     int
	assignedBy int
	assignedAt time.Time
}

func NewRoleAssignment(
	roleID int,
	userID int,
	assignedBy int,
) (RoleAssignment, error) {
	if roleID <= 0 {
		return RoleAssignment{}, NewError(CodeInvalidInput, "role ID must be positive")
	}

	if userID <= 0 {
		return RoleAssignment{}, NewError(CodeInvalidInput, "user ID must be positive")
	}

	if assignedBy <= 0 {
		return RoleAssignment{}, NewError(CodeInvalidInput, "assigned by ID must be positive")
	}

	return RoleAssignment{
		roleID:     roleID,
		userID:     userID,
		assignedBy: assignedBy,
		assignedAt: time.Now().UTC(),
	}, nil
}

func (a RoleAssignment) RoleID() int {
	return a.roleID
}

func (a RoleAssignment) UserID() int {
	return a.userID
}

func (a RoleAssignment) AssignedBy() int {
	return a.assignedBy
}

func (a RoleAssignment) AssignedAt() time.Time {
	return a.assignedAt
}
