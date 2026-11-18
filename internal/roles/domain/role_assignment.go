package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RoleAssignment struct {
	roleID     uuid.UUID
	userID     uuid.UUID
	assignedBy uuid.UUID
	assignedAt time.Time
}

func NewRoleAssignment(
	roleID uuid.UUID,
	userID uuid.UUID,
	assignedBy uuid.UUID,
) (RoleAssignment, error) {
	if roleID == uuid.Nil {
		return RoleAssignment{}, errors.New("role ID is empty")
	}

	if userID == uuid.Nil {
		return RoleAssignment{}, errors.New("user ID is empty")
	}

	if assignedBy == uuid.Nil {
		return RoleAssignment{}, errors.New("assigned by ID is empty")
	}

	return RoleAssignment{
		roleID:     roleID,
		userID:     userID,
		assignedBy: assignedBy,
		assignedAt: time.Now().UTC(),
	}, nil
}

func (a RoleAssignment) RoleID() uuid.UUID {
	return a.roleID
}

func (a RoleAssignment) UserID() uuid.UUID {
	return a.userID
}

func (a RoleAssignment) AssignedBy() uuid.UUID {
	return a.assignedBy
}

func (a RoleAssignment) AssignedAt() time.Time {
	return a.assignedAt
}
