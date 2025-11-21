package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"
)

type Order struct {
	OrderUID        string    `json:"order_uid" db:"order_uid" validate:"required, uuid"`
	TrackNumber     string    `json:"track_number" db:"track_number" validate:"required, min=1,max=50"`
	Entry           string    `json:"entry" db:"entry" validate:"required, min=1,max=50"`
	Locale          string    `json:"locale" db:"locale" validate:"required, min=1,max=50"`
	InternalSig     string    `json:"internal_signature" db:"internal_signature" validate:"required, min=1,max=50"`
	CustomerID      string    `json:"customer_id" db:"customer_id" validate:"required, min=1,max=50"`
	DeliveryService string    `json:"delivery_service" db:"delivery_service" validate:"required, min=1,max=50"`
	ShardKey        string    `json:"shardkey" db:"shardkey" validate:"required, min=1,max=50"`
	SmID            int       `json:"sm_id" db:"sm_id" validate:"required"`
	DateCreated     time.Time `json:"date_created" db:"date_created" validate:"required"`
	OofShard        string    `json:"oof_shard" db:"oof_shard" validate:"required, min=1,max=50"`
	Delivery        Delivery  `json:"delivery" db:"-"`
	Payment         Payment   `json:"payment" db:"-"`
	Items           []Item    `json:"items" db:"-"`
}

func (o *Order) Validate() error {
	if err := o.validateOrderUID(); err != nil {
		return err
	}
	if err := o.validateTrackNumber(); err != nil {
		return err
	}
	if err := o.validateEntry(); err != nil {
		return err
	}
	if err := o.validateLocale(); err != nil {
		return err
	}
	if err := o.validateCustomerID(); err != nil {
		return err
	}
	if err := o.validateDeliveryService(); err != nil {
		return err
	}
	if err := o.validateShardKey(); err != nil {
		return err
	}
	if err := o.validateSmID(); err != nil {
		return err
	}
	if err := o.validateDateCreated(); err != nil {
		return err
	}
	if err := o.validateOofShard(); err != nil {
		return err
	}
	if err := o.validateDelivery(); err != nil {
		return err
	}
	if err := o.validatePayment(); err != nil {
		return err
	}
	if err := o.validateItems(); err != nil {
		return err
	}

	return nil
}

func (o *Order) validateOrderUID() error {
	if o.OrderUID == "" {
		return errors.New("order_uid is required")
	}
	if utf8.RuneCountInString(o.OrderUID) > 50 {
		return errors.New("order_uid must be less than 50 characters")
	}
	return nil
}

func (o *Order) validateTrackNumber() error {
	if o.TrackNumber == "" {
		return errors.New("track_number is required")
	}
	if utf8.RuneCountInString(o.TrackNumber) > 100 {
		return errors.New("track_number must be less than 100 characters")
	}
	return nil
}

func (o *Order) validateEntry() error {
	if o.Entry == "" {
		return errors.New("entry is required")
	}
	if utf8.RuneCountInString(o.Entry) > 50 {
		return errors.New("entry must be less than 50 characters")
	}
	return nil
}

func (o *Order) validateLocale() error {
	if o.Locale == "" {
		return errors.New("locale is required")
	}
	// Проверяем, что locale соответствует формату (например, "ru", "en-US")
	matched, _ := regexp.MatchString(`^[a-z]{2}(-[A-Z]{2})?$`, o.Locale)
	if !matched {
		return errors.New("locale must be in format like 'ru' or 'en-US'")
	}
	return nil
}

func (o *Order) validateCustomerID() error {
	if o.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if utf8.RuneCountInString(o.CustomerID) > 50 {
		return errors.New("customer_id must be less than 50 characters")
	}
	return nil
}

func (o *Order) validateDeliveryService() error {
	if o.DeliveryService == "" {
		return errors.New("delivery_service is required")
	}
	if utf8.RuneCountInString(o.DeliveryService) > 100 {
		return errors.New("delivery_service must be less than 100 characters")
	}
	return nil
}

func (o *Order) validateShardKey() error {
	if o.ShardKey == "" {
		return errors.New("shardkey is required")
	}
	if utf8.RuneCountInString(o.ShardKey) > 50 {
		return errors.New("shardkey must be less than 50 characters")
	}
	return nil
}

func (o *Order) validateSmID() error {
	if o.SmID <= 0 {
		return errors.New("sm_id must be positive")
	}
	return nil
}

func (o *Order) validateDateCreated() error {
	if o.DateCreated.IsZero() {
		return errors.New("date_created is required")
	}
	if o.DateCreated.After(time.Now()) {
		return errors.New("date_created cannot be in the future")
	}
	return nil
}

func (o *Order) validateOofShard() error {
	if o.OofShard == "" {
		return errors.New("oof_shard is required")
	}
	if utf8.RuneCountInString(o.OofShard) > 50 {
		return errors.New("oof_shard must be less than 50 characters")
	}
	return nil
}

func (o *Order) validateDelivery() error {
	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("delivery validation failed: %v", err)
	}
	return nil
}

func (o *Order) validatePayment() error {
	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("payment validation failed: %v", err)
	}
	return nil
}

func (o *Order) validateItems() error {
	if len(o.Items) == 0 {
		return errors.New("at least one item is required")
	}
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d validation failed: %v", i, err)
		}
	}
	return nil
}
