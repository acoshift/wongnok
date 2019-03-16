package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/acoshift/wongnok/internal/auth"
)

// API handler
type API struct {
	Auth *auth.Auth
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
