package models

type Delivery struct {
	ID       int    `json:"id" db:"id"`
	OrderUID string `json:"order_uid" db:"order_uid" validate:"required, uuid"`
	Name     string `json:"name" db:"name" validate:"required, min=1,max=50"`
	Phone    string `json:"phone" db:"phone" validate:"required, e164"`
	Zip      string `json:"zip" db:"zip" validate:"required, min=1,max=50"`
	City     string `json:"city" db:"city" validate:"required, min=1,max=50"`
	Address  string `json:"address" db:"address" validate:"required, min=1,max=50"`
	Region   string `json:"region" db:"region" validate:"required, min=1,max=50"`
	Email    string `json:"email" db:"email" validate:"required, email"`
}
