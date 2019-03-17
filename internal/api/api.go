package api

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// API handler
type API struct {
	Auth AuthService
}

// AuthService type
type AuthService interface {
	SignUp(ctx context.Context, username, password string) (userID int64, err error)
	SignIn(ctx context.Context, username, password string) (token string, err error)
	SignOut(ctx context.Context, token string) error
}

// Handler returns api's handler
func (api API) Handler() http.Handler {
	router := httprouter.New()

	// auth
	router.POST("/auth/signup", api.authSignUp)
	router.POST("/auth/signin", api.authSignIn)
	router.POST("/auth/signout", api.authSignOut)

	return router
}

func decodeJSON(r *http.Request, v interface{}) error {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		return fmt.Errorf("invalid content-type")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}
