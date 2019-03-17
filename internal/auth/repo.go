package auth

import (
	"context"
	"database/sql"
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
	return
}

func (repo) DeleteToken(ctx context.Context, db *sql.DB, token string) error {
	_, err := db.ExecContext(ctx, `
		delete from auth_tokens
		where id = $1
	`, token)
	return err
}
