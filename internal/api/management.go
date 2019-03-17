package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/acoshift/wongnok/internal/management"
	"github.com/acoshift/wongnok/internal/validate"
)

func (api *API) managementCreateShop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Photos      []string `json:"photos"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	shopID, err := api.Management.CreateShop(ctx, &management.CreateShop{
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos,
	})
	if err, ok := err.(*validate.Error); ok {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	encodeJSON(w, struct {
		ID int64 `json:"id"`
	}{shopID})
}

func (api *API) managementListShops(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	shops, err := api.Management.ListShops(ctx)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	type item struct {
		ID          int64    `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Photos      []string `json:"photos"`
		CreatedAt   string   `json:"createdAt"`
	}
	list := make([]*item, 0, len(shops))
	for _, x := range shops {
		list = append(list, &item{
			ID:          x.ID,
			Name:        x.Name,
			Description: x.Description,
			Photos:      x.Photos,
			CreatedAt:   formatTime(x.CreatedAt),
		})
	}

	encodeJSON(w, list)
}
