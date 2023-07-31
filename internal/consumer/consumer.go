package consumer

import (
	"encoding/json"
	"fmt"
	"wb_test_1/internal/logger"
	"wb_test_1/internal/models/order"
	"wb_test_1/internal/order_service"

	"github.com/nats-io/stan.go"
)

type ConsumerServicer interface {
	Insert(order.Order) error
}

const (
	clusterID      = "test-cluster"
	clientID       = "consumer"
	DefaultNatsURL = "nats://127.0.0.1:4222"
)

func Start(orderService ConsumerServicer, c order_service.CacheRepositoryWorker, logger logger.Logger) {
	natsUrl := stan.NatsURL(DefaultNatsURL)
	sc, err := stan.Connect(clusterID, clientID, natsUrl)
	if err != nil {
		err := fmt.Errorf("connection to nats server: %v", err)
		logger.Error(err)
		return
	}

	defer sc.Close()

	sub, _ := sc.Subscribe("orders", func(m *stan.Msg) {
		var data order.Order
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			err = fmt.Errorf("order unmarshaling error: %v, data: %v", err, m.Data)
			logger.Error(err)
			return
		}
		err = orderService.Insert(data)
		if err != nil {
			err = fmt.Errorf("inserting order to db: %v", err)
			logger.Error(err)
			return
		}
		err = c.Insert(data.OrderUid, data)
		if err != nil {
			err = fmt.Errorf("inserting order to cache: %v", err)
			logger.Error(err)
			return
		}

		logger.Infof("Inserted order with id = %s", data.OrderUid)
	})

	defer sub.Unsubscribe()

	select {}
}
