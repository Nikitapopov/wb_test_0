package order

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

type OrderDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type OrderDeliveryValidation struct {
	Name    *string `json:"name" validate:"required"`
	Phone   *string `json:"phone" validate:"required"`
	Zip     *string `json:"zip" validate:"required"`
	City    *string `json:"city" validate:"required"`
	Address *string `json:"address" validate:"required"`
	Region  *string `json:"region" validate:"required"`
	Email   *string `json:"email" validate:"required"`
}

func (d *OrderDelivery) UnmarshalJSON(data []byte) (err error) {
	all := OrderDeliveryValidation{}
	dec := json.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&all); err != nil {
		return err
	}
	*d = OrderDelivery{
		Name:    *all.Name,
		Phone:   *all.Phone,
		Zip:     *all.Zip,
		City:    *all.City,
		Address: *all.Address,
		Region:  *all.Region,
		Email:   *all.Email,
	}

	validate := validator.New()
	err = validate.Struct(d)
	if err != nil {
		return err
	}

	return
}

func (d *OrderDelivery) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}

func (d OrderDelivery) Value() (driver.Value, error) {
	return json.Marshal(d)
}

type OrderPayment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type OrderPaymentValidation struct {
	Transaction  *string `json:"transaction" validate:"required"`
	RequestId    *string `json:"request_id" validate:"required"`
	Currency     *string `json:"currency" validate:"required"`
	Provider     *string `json:"provider" validate:"required"`
	Amount       *int    `json:"amount" validate:"required"`
	PaymentDt    *int    `json:"payment_dt" validate:"required"`
	Bank         *string `json:"bank" validate:"required"`
	DeliveryCost *int    `json:"delivery_cost" validate:"required"`
	GoodsTotal   *int    `json:"goods_total" validate:"required"`
	CustomFee    *int    `json:"custom_fee" validate:"required"`
}

func (p *OrderPayment) UnmarshalJSON(data []byte) (err error) {
	all := OrderPaymentValidation{}
	dec := json.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&all); err != nil {
		return err
	}
	*p = OrderPayment{
		Transaction:  *all.Transaction,
		RequestId:    *all.RequestId,
		Currency:     *all.Currency,
		Provider:     *all.Provider,
		Amount:       *all.Amount,
		PaymentDt:    *all.PaymentDt,
		Bank:         *all.Bank,
		DeliveryCost: *all.DeliveryCost,
		GoodsTotal:   *all.GoodsTotal,
		CustomFee:    *all.CustomFee,
	}

	validate := validator.New()
	err = validate.Struct(p)
	if err != nil {
		return err
	}

	return
}

func (p *OrderPayment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

func (p OrderPayment) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type OrderItem struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type OrderItemValidation struct {
	ChrtId      *int    `json:"chrt_id" validate:"required"`
	TrackNumber *string `json:"track_number" validate:"required"`
	Price       *int    `json:"price" validate:"required"`
	Rid         *string `json:"rid" validate:"required"`
	Name        *string `json:"name" validate:"required"`
	Sale        *int    `json:"sale" validate:"required"`
	Size        *string `json:"size" validate:"required"`
	TotalPrice  *int    `json:"total_price" validate:"required"`
	NmId        *int    `json:"nm_id" validate:"required"`
	Brand       *string `json:"brand" validate:"required"`
	Status      *int    `json:"status" validate:"required"`
}

func (i *OrderItem) UnmarshalJSON(data []byte) (err error) {
	all := OrderItemValidation{}
	dec := json.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&all); err != nil {
		return err
	}
	*i = OrderItem{
		ChrtId:      *all.ChrtId,
		TrackNumber: *all.TrackNumber,
		Price:       *all.Price,
		Rid:         *all.Rid,
		Name:        *all.Name,
		Sale:        *all.Sale,
		Size:        *all.Size,
		TotalPrice:  *all.TotalPrice,
		NmId:        *all.NmId,
		Brand:       *all.Brand,
		Status:      *all.Status,
	}

	validate := validator.New()
	err = validate.Struct(i)
	if err != nil {
		return err
	}

	return
}

type OrderItems []OrderItem

func (o *OrderItems) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, o)
	case string:
		return json.Unmarshal([]byte(v), o)
	}
	return errors.New("type assertion failed")
}

func (o OrderItems) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type Order struct {
	OrderUid          string        `json:"order_uid"`
	TrackNumber       string        `json:"track_number"`
	Entry             string        `json:"entry"`
	Delivery          OrderDelivery `json:"delivery"`
	Payment           OrderPayment  `json:"payment"`
	Items             OrderItems    `json:"items"`
	Locale            string        `json:"locale"`
	InternalSignature string        `json:"internal_signature"`
	CustomerId        string        `json:"customer_id"`
	DeliveryService   string        `json:"delivery_service"`
	Shardkey          string        `json:"shardkey"`
	SmId              int           `json:"sm_id"`
	DateCreated       string        `json:"date_created"`
	OofShard          string        `json:"oof_shard"`
}

type OrderValidation struct {
	OrderUid          *string        `json:"order_uid" validate:"required"`
	TrackNumber       *string        `json:"track_number" validate:"required"`
	Entry             *string        `json:"entry" validate:"required"`
	Delivery          *OrderDelivery `json:"delivery" validate:"required"`
	Payment           *OrderPayment  `json:"payment" validate:"required"`
	Items             *OrderItems    `json:"items" validate:"required"`
	Locale            *string        `json:"locale" validate:"required"`
	InternalSignature *string        `json:"internal_signature" validate:"required"`
	CustomerId        *string        `json:"customer_id" validate:"required"`
	DeliveryService   *string        `json:"delivery_service" validate:"required"`
	Shardkey          *string        `json:"shardkey" validate:"required"`
	SmId              *int           `json:"sm_id" validate:"required"`
	DateCreated       *string        `json:"date_created" validate:"required"`
	OofShard          *string        `json:"oof_shard" validate:"required"`
}

func (o *Order) UnmarshalJSON(data []byte) (err error) {
	all := OrderValidation{}
	dec := json.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&all); err != nil {
		return err
	}
	*o = Order{
		OrderUid:          *all.OrderUid,
		TrackNumber:       *all.TrackNumber,
		Entry:             *all.Entry,
		Delivery:          *all.Delivery,
		Payment:           *all.Payment,
		Items:             *all.Items,
		Locale:            *all.Locale,
		InternalSignature: *all.InternalSignature,
		CustomerId:        *all.CustomerId,
		DeliveryService:   *all.DeliveryService,
		Shardkey:          *all.Shardkey,
		SmId:              *all.SmId,
		DateCreated:       *all.DateCreated,
		OofShard:          *all.OofShard,
	}

	validate := validator.New()
	err = validate.Struct(o)
	if err != nil {
		return err
	}

	return
}

func (o *Order) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &o)
}

func (o Order) Value() (driver.Value, error) {
	return json.Marshal(o)
}
