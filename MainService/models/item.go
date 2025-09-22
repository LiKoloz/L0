package models

import "errors"

type Item struct {
	ID          int    `json:"id" db:"id"`
	OrderUID    string `json:"order_uid" db:"order_uid"`
	ChrtID      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	RID         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmID        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
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
