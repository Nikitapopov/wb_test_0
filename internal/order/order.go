package order

import (
	"fmt"
	"wb_test_1/internal/logger"
)

type DbRepositoryWorker interface {
	GetById(id string) (*Order, error)
	GetAll() ([]*Order, error)
	Insert(Order) error
}

type CacheRepositoryWorker interface {
	GetById(string) (*Order, error)
	GetIdsList() []string
	Insert(string, Order) error
}

type Service struct {
	dbRepo    DbRepositoryWorker
	cacheRepo CacheRepositoryWorker
	logger    logger.Logger
}

func NewService(dbRepo DbRepositoryWorker, cacheRepo CacheRepositoryWorker, logger logger.Logger) OrderServicer {
	service := &Service{
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

func (s *Service) syncCacheWithDb() error {
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

func (s *Service) GetById(id string) (*Order, error) {
	return s.cacheRepo.GetById(id)
}

func (s *Service) GetIdsList() []string {
	return s.cacheRepo.GetIdsList()
}

func (s *Service) Insert(order Order) error {
	if err := s.dbRepo.Insert(order); err != nil {
		return fmt.Errorf("inserting order to db: %v", err)
	}

	if err := s.cacheRepo.Insert(order.OrderUid, order); err != nil {
		return fmt.Errorf("inserting order to cache: %v", err)
	}

	return nil
}
