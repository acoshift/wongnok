package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (api *API) authSignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	_, err = api.Auth.SignUp(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, struct {
		Success bool `json:"success"`
	}{true})
}

func (api *API) authSignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	token, err := api.Auth.SignIn(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}{true, token})
}

func (api *API) authSignOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Token string `json:"token"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = api.Auth.SignOut(ctx, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, struct {
		Success bool `json:"success"`
	}{true})
}
