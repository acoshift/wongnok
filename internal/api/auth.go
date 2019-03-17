package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/acoshift/wongnok/internal/auth"
	"github.com/acoshift/wongnok/internal/validate"
)

func (api *API) authSignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	_, err = api.Auth.SignUp(ctx, req.Username, req.Password)

	// case 1: many error values
	if err == auth.ErrUsernameRequired {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	if err == auth.ErrUsernameTooShort {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	if err == auth.ErrUsernameTooLong {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	if err == auth.ErrUsernameInvalid {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	if err == auth.ErrUsernameNotAvailable {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	// case 2: group error using type
	if err, ok := err.(*validate.Error); ok {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
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
