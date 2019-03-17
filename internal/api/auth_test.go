package api

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestAPI_authSignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var called bool
		api := API{Auth: &mockAuthSignUp{
			Func: func(ctx context.Context, username, password string) (userID int64, err error) {
				called = true
				assert.Equal(t, "tester", username)
				assert.Equal(t, "123456", password)
				return 10, nil
			},
		}}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/",
			strings.NewReader( /* language=JSON */ `
			{
				"username": "tester",
				"password": "123456"
			}
		`))
		r.Header.Set("Content-Type", "application/json; charset=utf-8")

		api.authSignUp(w, r, httprouter.Params{})

		assert.True(t, called)
		assert.EqualValues(t, 200, w.Code)
		// language=JSON
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}

type mockAuthSignUp struct {
	AuthService

	Func func(ctx context.Context, username, password string) (userID int64, err error)
}

func (m *mockAuthSignUp) SignUp(ctx context.Context, username, password string) (userID int64, err error) {
	return m.Func(ctx, username, password)
}
