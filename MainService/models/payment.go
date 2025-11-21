package models

import "errors"

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction" validate:"required, min=1,max=50"`
	RequestID    string `json:"request_id" db:"request_id" validate:"required, min=1,max=50"`
	Currency     string `json:"currency" db:"currency" validate:"required, min=1,max=50"`
	Provider     string `json:"provider" db:"provider" validate:"required, min=1,max=50"`
	Amount       int    `json:"amount" db:"amount" validate:"required"`
	PaymentDT    int64  `json:"payment_dt" db:"payment_dt" validate:"required"`
	Bank         string `json:"bank" db:"bank" validate:"required, min=1,max=50"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"required"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"required"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"required"`
}

func (p *Payment) Validate() error {
	if p.Transaction == "" {
		return errors.New("payment transaction is required")
	}

	if p.Currency == "" {
		return errors.New("currency is required")
	}

	if p.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if p.PaymentDT == 0 {
		return errors.New("payment datetime is required")
	}

	return nil
}
