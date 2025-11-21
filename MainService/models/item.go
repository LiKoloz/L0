package models

import "errors"

type Item struct {
	ID          int    `json:"id" db:"id"`
	OrderUID    string `json:"order_uid" db:"order_uid" validate:"required, uuid"`
	ChrtID      int    `json:"chrt_id" db:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required, min=1,max=50"`
	Price       int    `json:"price" db:"price" validate:"required, min=1"`
	RID         string `json:"rid" db:"rid" validate:"required, min=1,max=50"`
	Name        string `json:"name" db:"name" validate:"required, min=1,max=50"`
	Sale        int    `json:"sale" db:"sale" validate:"required, min=1"`
	Size        string `json:"size" db:"size" validate:"required, min=1,max=50"`
	TotalPrice  int    `json:"total_price" db:"total_price" validate:"required, min=1"`
	NmID        int    `json:"nm_id" db:"nm_id" validate:"required"`
	Brand       string `json:"brand" db:"brand" validate:"required, min=1,max=50"`
	Status      int    `json:"status" db:"status" validate:"required"`
}

func (i *Item) Validate() error {
	if i.ChrtID <= 0 {
		return errors.New("chrt_id must be positive")
	}

	if i.Price <= 0 {
		return errors.New("price must be positive")
	}

	if i.Name == "" {
		return errors.New("item name is required")
	}

	if i.TotalPrice <= 0 {
		return errors.New("total_price must be positive")
	}

	return nil
}
