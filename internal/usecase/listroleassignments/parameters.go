package listroleassignments

import (
	"time"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	"github.com/hedgehog-pm/roles-service/internal/pagination"
)

type (
	Params struct {
		RoleID    int
		UserID    int
		ScopeType domain.ScopeType
		ScopeID   int
		pagination.Pagination
	}
	Assignment struct {
		RoleID     int
		UserID     int
		RoleName   string
		Scope      domain.ScopeRef
		AssignedBy int
		AssignedAt time.Time
	}
)
