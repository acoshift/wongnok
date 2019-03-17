package management

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"

	"github.com/acoshift/wongnok/internal/validate"
)

// Management service
type Management struct {
	db *sql.DB
}

// New creates new management service
func New(db *sql.DB) *Management {
	return &Management{db}
}

// CreateShop type
type CreateShop struct {
	Name        string
	Description string
	Photos      []string
}

// CreateShop creates new shop
func (svc *Management) CreateShop(ctx context.Context, shop *CreateShop) (shopID int64, err error) {
	if shop.Name == "" {
		return 0, validate.NewRequiredError("name")
	}
	if utf8.RuneCountInString(shop.Name) > 100 {
		return 0, validate.NewError("name", "too long")
	}
	if shop.Description == "" {
		return 0, validate.NewRequiredError("description")
	}
	if utf8.RuneCountInString(shop.Description) > 2000 {
		return 0, validate.NewError("description", "too long")
	}
	if len(shop.Photos) > 10 {
		return 0, validate.NewError("photos", "limit to 10 photos")
	}
	for i, photo := range shop.Photos {
		if l := len(photo); l == 0 {
			return 0, validate.NewError(
				fmt.Sprintf("photos[%d]", i),
				"photo url empty",
			)
		} else if l > 200 {
			return 0, validate.NewError(
				fmt.Sprintf("photos[%d]", i),
				"photo url too long",
			)
		}
		if !govalidator.IsURL(photo) {
			return 0, validate.NewError(
				fmt.Sprintf("photos[%d]", i),
				"photo is not an url",
			)
		}
	}

	err = svc.db.QueryRowContext(ctx, `
		insert into shops
			(name, description, photos)
		values
			($1, $2, $3)
		returning id
	`, shop.Name, shop.Description, pq.Array(shop.Photos)).Scan(&shopID)
	if err != nil {
		return 0, err
	}
	return shopID, nil
}

// Shop entity
type Shop struct {
	ID          int64
	Name        string
	Description string
	Photos      []string
	CreatedAt   time.Time
}

// ListShops retrieves all shops
func (svc *Management) ListShops(ctx context.Context) ([]*Shop, error) {
	rows, err := svc.db.QueryContext(ctx, `
		select
			id, name, description, photos, created_at
		from shops
		order by id desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []*Shop
	for rows.Next() {
		var shop Shop
		err = rows.Scan(
			&shop.ID, &shop.Name, &shop.Description,
			pq.Array(&shop.Photos), &shop.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shops = append(shops, &shop)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return shops, nil
}
