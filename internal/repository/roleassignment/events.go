package roleassignment

import (
	"context"
	"encoding/json"

	"github.com/stephenafamo/bob/types"

	rolesdb "github.com/hedgehog-pm/roles-service/internal/generated/bob/roleassignment"
	"github.com/hedgehog-pm/roles-service/internal/storage/postgres"
)

func event(ctx context.Context, tx executor, kind string, id int, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = postgres.QueryExec(ctx, tx, rolesdb.InsertOutboxMessage([]rolesdb.InsertOutboxMessage_Rows{{
		EventType: kind, AggregateID: id, Payload: types.NewJSON(json.RawMessage(data)),
	}}))

	return err
}
