package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/acoshift/wongnok/internal/management"
)

// API handler
type API struct {
	Auth       AuthService
	Management *management.Management
}

// AuthService type
type AuthService interface {
	SignUp(ctx context.Context, username, password string) (userID int64, err error)
	SignIn(ctx context.Context, username, password string) (token string, err error)
	SignOut(ctx context.Context, token string) error
	VerifyToken(ctx context.Context, token string) (userID int64, isAdmin bool, err error)
}

// Handler returns api's handler
func (api API) Handler() http.Handler {
	router := httprouter.New()

	router.GET("/healthz", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	// auth
	router.POST("/auth/signup", api.authSignUp)
	router.POST("/auth/signin", api.authSignIn)
	router.POST("/auth/signout", api.authSignOut)

	// management
	{
		router := newGroupRouter(router, "/management", onlyAdminGuard)
		router.POST("/shops", api.managementCreateShop)
		router.GET("/shops", api.managementListShops)
	}

	return api.fetchCredential(router)
}

type groupRouter struct {
	router     *httprouter.Router
	prefix     string
	middleware func(httprouter.Handle) httprouter.Handle
}

func (router *groupRouter) GET(path string, h httprouter.Handle) {
	router.router.GET(router.prefix+path, router.middleware(h))
}

func (router *groupRouter) POST(path string, h httprouter.Handle) {
	router.router.POST(router.prefix+path, router.middleware(h))
}

func newGroupRouter(router *httprouter.Router, prefix string, middleware func(httprouter.Handle) httprouter.Handle) *groupRouter {
	return &groupRouter{
		router,
		prefix,
		middleware,
	}
}

type ctxKey string

const (
	ctxKeyUserID  ctxKey = "user_id"
	ctxKeyIsAdmin ctxKey = "is_admin"
)

func (api *API) fetchCredential(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			h.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		userID, isAdmin, err := api.Auth.VerifyToken(ctx, token)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		ctx = context.WithValue(ctx, ctxKeyUserID, userID)
		ctx = context.WithValue(ctx, ctxKeyIsAdmin, isAdmin)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func getUserID(ctx context.Context) int64 {
	x, _ := ctx.Value(ctxKeyUserID).(int64)
	return x
}

func getIsAdmin(ctx context.Context) bool {
	x, _ := ctx.Value(ctxKeyIsAdmin).(bool)
	return x
}

func onlyAdminGuard(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		if !getIsAdmin(ctx) {
			handleError(w, http.StatusForbidden, fmt.Errorf("forbidden"))
			return
		}
		h(w, r, ps)
	}
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

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if statusCode == http.StatusInternalServerError {
		log.Println(err)
		err = fmt.Errorf("internal error")
	}
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
