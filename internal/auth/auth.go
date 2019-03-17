package auth

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/acoshift/wongnok/internal/validate"
)

// Auth service
type Auth struct {
	db   *sql.DB
	repo repository
}

type repository interface {
	InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error)
	DeleteToken(ctx context.Context, db *sql.DB, token string) error
}

// New creates new auth service
func New(db *sql.DB) *Auth {
	return &Auth{db, repo{}}
}

var reUsername = regexp.MustCompile(`^[a-z0-9]*$`)

// SignUp registers new user
func (svc *Auth) SignUp(ctx context.Context, username, password string) (userID int64, err error) {
	// normalize data
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)

	// validate
	if username == "" {
		return 0, ErrUsernameRequired
	}
	if len(username) < 4 {
		return 0, ErrUsernameTooShort
	}
	if len(username) > 20 {
		return 0, ErrUsernameTooLong
	}
	if !reUsername.MatchString(username) {
		return 0, ErrUsernameInvalid
	}
	if password == "" {
		return 0, validate.NewRequiredError("password")
	}
	if len(password) < 6 {
		return 0, validate.NewError("password", "too short")
	}
	if len(password) > 64 {
		return 0, validate.NewError("password", "too long")
	}

	// hash password
	hashedPass := hashPassword(password)

	userID, err = svc.repo.InsertUser(ctx, svc.db, username, hashedPass)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// SignIn sign in user
func (svc *Auth) SignIn(ctx context.Context, username, password string) (token string, err error) {
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)

	if username == "" {
		return "", fmt.Errorf("username required")
	}
	if len(username) > 20 {
		return "", fmt.Errorf("username too long")
	}
	if password == "" {
		return "", fmt.Errorf("password required")
	}
	if len(password) > 64 {
		return "", fmt.Errorf("password too long")
	}

	var (
		userID       int64
		userPassword string
	)
	err = svc.db.QueryRowContext(ctx, `
		select
			id, password
		from users
		where username = $1
	`, username).Scan(&userID, &userPassword)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	if !compareHashAndPassword(userPassword, password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token = generateToken()

	_, err = svc.db.ExecContext(ctx, `
		insert into auth_tokens
			(id, user_id)
		values
			($1, $2)
	`, token, userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// SignOut sign out user
func (svc *Auth) SignOut(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token required")
	}

	return svc.repo.DeleteToken(ctx, svc.db, token)
}

// VerifyToken returns user id if token valid
func (svc *Auth) VerifyToken(ctx context.Context, token string) (userID int64, isAdmin bool, err error) {
	if token == "" {
		return 0, false, nil
	}

	err = svc.db.QueryRowContext(ctx, `
		select
			auth_tokens.user_id, users.is_admin
		from auth_tokens
		left join users on auth_tokens.user_id = users.id
		where auth_tokens.id = $1
	`, token).Scan(&userID, &isAdmin)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return userID, isAdmin, nil
}
