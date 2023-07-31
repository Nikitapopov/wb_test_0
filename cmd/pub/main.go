package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"wb_test_1/internal/models/order"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	stan "github.com/nats-io/stan.go"
)

const (
	clusterID      = "test-cluster"
	clientID       = "publisher"
	DefaultNatsURL = "nats://127.0.0.1:4222"
)

func main() {
	natsUrl := stan.NatsURL(DefaultNatsURL)
	sc, err := stan.Connect(clusterID, clientID, natsUrl)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	for {
		data, err := generateOrder()
		if err != nil {
			fmt.Printf("Generating order: %v", err)
			continue
		}
		dataBytes, _ := json.Marshal(data)
		sc.Publish("orders", dataBytes)

		// Генерация заявок через промежуток от 0 до 10 секунд
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func generateOrder() (*order.Order, error) {
	address := faker.GetRealAddress()
	trackNumber := faker.UUIDDigit()
	date, err := time.Parse("2006-01-02 15:04:05", faker.Timestamp())
	if err != nil {
		return nil, fmt.Errorf("generating fake timestamp: %w", err)
	}
	formattedDate := date.Format("2006-01-02T15:04:05Z")

	itemsNumber := rand.Intn(3)
	items := make(order.OrderItems, 0, 1)
	for i := 0; i < itemsNumber; i++ {
		items = append(items, order.OrderItem{
			ChrtId:      rand.Intn(9000000) + 1000000,
			TrackNumber: trackNumber,
			Price:       rand.Intn(10000),
			Rid:         faker.UUIDDigit(),
			Name:        faker.Name(),
			Sale:        rand.Intn(100),
			Size:        fmt.Sprint(rand.Intn(10)),
			TotalPrice:  rand.Intn(10000),
			NmId:        rand.Intn(9000000) + 1000000,
			Brand:       fmt.Sprintf("%s %s", faker.FirstName(), faker.LastName()),
			Status:      rand.Intn(500) + 100,
		})
	}

	return &order.Order{
		OrderUid:    faker.UUIDDigit(),
		TrackNumber: trackNumber,
		Entry:       faker.Word(options.WithRandomStringLength(4)),
		Delivery: order.OrderDelivery{
			Name:    fmt.Sprintf("%s %s", faker.FirstName(), faker.LastName()),
			Phone:   faker.Phonenumber(),
			Zip:     address.PostalCode,
			City:    address.City,
			Address: address.Address,
			Region:  address.State,
			Email:   faker.Email(),
		},
		Payment: order.OrderPayment{
			Transaction:  faker.UUIDDigit(),
			RequestId:    faker.UUIDDigit(),
			Currency:     faker.Currency(),
			Provider:     faker.Word(options.WithRandomStringLength(5)),
			Amount:       rand.Intn(10000),
			PaymentDt:    rand.Intn(9000000000) + 1000000000,
			Bank:         faker.Name(),
			DeliveryCost: rand.Intn(10000),
			GoodsTotal:   rand.Intn(10000),
			CustomFee:    rand.Intn(100),
		},
		Items:             items,
		Locale:            faker.Word(options.WithRandomStringLength(2)),
		InternalSignature: faker.Word(options.WithRandomStringLength(3)),
		CustomerId:        faker.UUIDDigit(),
		DeliveryService:   faker.Word(options.WithRandomStringLength(5)),
		Shardkey:          fmt.Sprint(rand.Intn(10)),
		SmId:              rand.Intn(9000000) + 1000000,
		DateCreated:       formattedDate,
		OofShard:          fmt.Sprint(rand.Intn(10)),
	}, nil
}
