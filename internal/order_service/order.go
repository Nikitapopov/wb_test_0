package order_service

import (
	"fmt"
	"wb_test_1/internal/logger"
	"wb_test_1/internal/models/order"
)

type DbRepositoryWorker interface {
	GetById(id string) (*order.Order, error)
	GetAll() ([]*order.Order, error)
	Insert(order.Order) error
}

type CacheRepositoryWorker interface {
	GetById(string) (*order.Order, error)
	GetIdsList() []string
	Insert(string, order.Order) error
}

type orderService struct {
	dbRepo    DbRepositoryWorker
	cacheRepo CacheRepositoryWorker
	logger    logger.Logger
}

func NewService(dbRepo DbRepositoryWorker, cacheRepo CacheRepositoryWorker, logger logger.Logger) OrderServicer {
	service := &orderService{
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
	err := service.syncCacheWithDb()
	if err != nil {
		err = fmt.Errorf("could not sync cache with db: %v", err)
		logger.Error(err)
	}
	return service
}

func (s *orderService) syncCacheWithDb() error {
	orders, err := s.dbRepo.GetAll()
	if err != nil {
		return fmt.Errorf("receiving all orders from db: %v", err)
	}
	for _, v := range orders {
		err = s.cacheRepo.Insert(v.OrderUid, *v)
		if err != nil {
			return fmt.Errorf("adding order during sync with db: %v", err)
		}
	}
	return nil
}

func (s *orderService) GetById(id string) (*order.Order, error) {
	return s.cacheRepo.GetById(id)
}

func (s *orderService) GetIdsList() []string {
	return s.cacheRepo.GetIdsList()
}

func (s *orderService) Insert(order order.Order) error {
	if err := s.dbRepo.Insert(order); err != nil {
		return fmt.Errorf("inserting order to db: %v", err)
	}

	if err := s.cacheRepo.Insert(order.OrderUid, order); err != nil {
		return fmt.Errorf("inserting order to cache: %v", err)
	}

	return nil
}
