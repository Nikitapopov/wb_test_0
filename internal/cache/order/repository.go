package order_repo_cache

import (
	"errors"
	"fmt"
	"wb_test_1/internal/logger"
	"wb_test_1/internal/models/order"
	"wb_test_1/internal/order_service"
	"wb_test_1/pkg/go_cache"
)

type repository struct {
	client go_cache.GoCacher
	logger logger.Logger
}

func NewRepo(client go_cache.GoCacher, logger logger.Logger) order_service.CacheRepositoryWorker {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) GetById(id string) (*order.Order, error) {
	cacheItem, ok := r.client.Get(id)
	if !ok {
		return nil, errors.New("order not found")
	}

	cacheOrder, ok := cacheItem.(order.Order)
	if !ok {
		return nil, errors.New("type assertion to order is failed")
	}

	return &cacheOrder, nil
}

func (r *repository) GetIdsList() []string {
	return r.client.GetKeys()
}

func (r *repository) Insert(id string, order order.Order) error {
	if err := r.client.Add(id, order); err != nil {
		err = fmt.Errorf("adding order to cache: %w, order object: %v", err, order)
		r.logger.Error(err)
		return err
	}
	return nil
}
