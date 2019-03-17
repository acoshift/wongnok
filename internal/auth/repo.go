package auth

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type repo struct{}

func (repo) InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
	err = db.QueryRowContext(ctx, `
		insert into users
			(username, password)
		values
			($1, $2)
		returning id
	`, username, password).Scan(&userID)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" && err.Constraint == "users_username_idx" {
			return 0, ErrUsernameNotAvailable
		}
	}
	return
}

func (repo) DeleteToken(ctx context.Context, db *sql.DB, token string) error {
	_, err := db.ExecContext(ctx, `
		delete from auth_tokens
		where id = $1
	`, token)
	return err
}
