package roletemplate

import "github.com/hedgehog-pm/roles-service/internal/domain"

type Filter struct {
	SystemID        int
	TargetScopeType domain.ScopeType
}
