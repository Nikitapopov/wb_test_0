package order_service

import "wb_test_1/internal/models/order"

type OrderServicer interface {
	GetById(id string) (*order.Order, error)
	GetIdsList() []string
	Insert(order.Order) error
}
