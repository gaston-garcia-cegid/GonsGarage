package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// fetchUserModelsByIDs loads active users (deleted_at IS NULL) by primary key.
func fetchUserModelsByIDs(ctx context.Context, x *sqlx.DB, ids []uuid.UUID) (map[uuid.UUID]UserModel, error) {
	if x == nil || len(ids) == 0 {
		return map[uuid.UUID]UserModel{}, nil
	}
	uniq := make([]uuid.UUID, 0, len(ids))
	seen := make(map[uuid.UUID]struct{})
	for _, id := range ids {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	if len(uniq) == 0 {
		return map[uuid.UUID]UserModel{}, nil
	}
	ph := make([]string, len(uniq))
	args := make([]interface{}, len(uniq))
	for i, id := range uniq {
		ph[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	q := `SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
FROM users WHERE deleted_at IS NULL AND id IN (` + strings.Join(ph, ",") + `)`
	var users []UserModel
	if err := x.SelectContext(ctx, &users, q, args...); err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]UserModel, len(users))
	for _, u := range users {
		out[u.ID] = u
	}
	return out, nil
}
