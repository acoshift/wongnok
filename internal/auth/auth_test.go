package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

var bgCtx = context.Background()

func Test_SignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := Auth{repo: &mockRepo{
			InsertUserFunc: func(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
				if username != "tester" {
					t.Errorf("expected username to be tester; got %s", username)
				}
				if !compareHashAndPassword(password, "123456") {
					t.Errorf("expected password to be equal to '123456'")
				}
				return 1, nil
			},
		}}
		userID, err := svc.SignUp(bgCtx, "Tester", "123456")
		if err != nil {
			t.Errorf("expected error to be nil; got %s", err)
		}
		if userID != 1 {
			t.Errorf("expected user id to be 1; got %d", userID)
		}
	})

	t.Run("DB Error", func(t *testing.T) {
		svc := Auth{repo: &mockRepoInsertUser{
			UserID: 0,
			Err:    fmt.Errorf("db error"),
		}}
		userID, err := svc.SignUp(bgCtx, "tester", "123456")
		if err == nil {
			t.Errorf("expected error return")
		}
		if userID != 0 {
			t.Errorf("expected user id to be 0; got %d", userID)
		}
	})

	t.Run("Username empty", func(t *testing.T) {
		svc := Auth{repo: &fakeRepo{}}
		userID, err := svc.SignUp(bgCtx, "", "123456")
		if err == nil {
			t.Errorf("expected error return")
		}
		if userID != 0 {
			t.Errorf("expected user id to be 0; got %d", userID)
		}
	})

	t.Run("Username too short", func(t *testing.T) {
		svc := Auth{repo: &fakeRepo{}}
		userID, err := svc.SignUp(bgCtx, "a", "123456")
		if err == nil {
			t.Errorf("expected error return")
		}
		if userID != 0 {
			t.Errorf("expected user id to be 0; got %d", userID)
		}
	})

	t.Run("Username too long", func(t *testing.T) {
		svc := Auth{repo: &fakeRepo{}}
		userID, err := svc.SignUp(bgCtx, strings.Repeat("a", 30), "123456")
		if err == nil {
			t.Errorf("expected error return")
		}
		if userID != 0 {
			t.Errorf("expected user id to be 0; got %d", userID)
		}
	})
}

type mockRepo struct {
	InsertUserFunc func(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error)
}

func (m *mockRepo) InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
	return m.InsertUserFunc(ctx, db, username, password)
}

type mockRepoInsertUser struct {
	repository

	UserID int64
	Err    error
}

func (m *mockRepoInsertUser) InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
	return m.UserID, m.Err
}

type fakeRepoItem struct {
	ID       int64
	Username string
	Password string
}

type fakeRepo struct {
	storage []*fakeRepoItem
}

func (f *fakeRepo) isUsernameExists(username string) bool {
	for _, it := range f.storage {
		if it.Username == username {
			return true
		}
	}
	return false
}

func (f *fakeRepo) InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
	if f.isUsernameExists(username) {
		return 0, fmt.Errorf("duplicate username")
	}

	id := int64(len(f.storage) + 1)

	f.storage = append(f.storage, &fakeRepoItem{ID: id, Username: username, Password: password})
	return id, nil
}
