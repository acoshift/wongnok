package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bgCtx = context.Background()

func TestAuth_SignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var called bool
		svc := Auth{repo: &mockRepo{
			InsertUserFunc: func(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
				called = true
				assert.Equal(t, "tester", username)
				assert.True(t, compareHashAndPassword(password, "123456"))
				return 1, nil
			},
		}}
		userID, err := svc.SignUp(bgCtx, "Tester", "123456")
		assert.NoError(t, err)
		assert.EqualValues(t, 1, userID)
		assert.True(t, called, "expected InsertUser was called")
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

func TestAuth_SignOut(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var called bool
		svc := Auth{repo: &mockRepo{
			DeleteTokenFunc: func(ctx context.Context, db *sql.DB, token string) error {
				called = true
				assert.Equal(t, "test-1234", token)
				return nil
			},
		}}
		err := svc.SignOut(bgCtx, "test-1234")
		assert.NoError(t, err)
		assert.True(t, called, "expected DeleteToken was called")
	})

	t.Run("Empty token", func(t *testing.T) {
		svc := Auth{repo: &mockRepo{
			DeleteTokenFunc: func(ctx context.Context, db *sql.DB, token string) error {
				assert.Fail(t, "expected DeleteToken was not be called")
				return nil
			},
		}}
		err := svc.SignOut(bgCtx, "")
		assert.Error(t, err)
	})
}

type mockRepo struct {
	InsertUserFunc  func(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error)
	DeleteTokenFunc func(ctx context.Context, db *sql.DB, token string) error
}

func (m *mockRepo) InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error) {
	return m.InsertUserFunc(ctx, db, username, password)
}

func (m *mockRepo) DeleteToken(ctx context.Context, db *sql.DB, token string) error {
	return m.DeleteTokenFunc(ctx, db, token)
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

func (f *fakeRepo) DeleteToken(ctx context.Context, db *sql.DB, token string) error {
	return nil
}
