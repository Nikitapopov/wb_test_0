package manager

import "wb_test_1/internal/models/order"

type ManagerServicer interface {
	GetById(id string) (*order.Order, error)
	GetKeys() []string
}

type Service struct {
	orderService ManagerServicer
}

func NewService(orderService ManagerServicer) Servicer {
	return &Service{orderService: orderService}
}

func (s *Service) GetById(id string) (*order.Order, error) {
	return s.orderService.GetById(id)
}

func (s *Service) GetIdsList() []string {
	return s.orderService.GetKeys()
}
