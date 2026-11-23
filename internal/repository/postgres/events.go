package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/types"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob"
)

type rolePayload struct {
	RoleID    int                     `json:"role_id"`
	ScopeType domain.ScopeType        `json:"scope_type"`
	ScopeID   int                     `json:"scope_id"`
	Version   int64                   `json:"version"`
	Codes     []domain.PermissionCode `json:"codes"`
}

func roleEvent(ctx context.Context, tx pgx.Tx, kind string, role *domain.Role) error {
	return event(ctx, tx, kind, role.ID, rolePayload{RoleID: role.ID, ScopeType: role.Scope.Type(), ScopeID: role.Scope.ID(), Version: role.Version, Codes: role.Permissions.Codes()})
}

func event(ctx context.Context, tx pgx.Tx, kind string, id int, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = queryExec(ctx, tx, rolesdb.InsertOutboxMessage([]rolesdb.InsertOutboxMessage_Rows{{
		EventType: kind, AggregateID: id, Payload: types.NewJSON(json.RawMessage(data)),
	}}))

	return err
}
