package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"wb_test_1/internal/order"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	stan "github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "publisher"
	DefaultNatsURL := "nats://127.0.0.1:4222"
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
		// dataBytes := []byte{123, 34, 111, 114, 100, 100, 101, 114, 95, 117, 105, 100, 34, 58, 34, 53, 48, 100, 53, 99, 99, 97, 52, 97, 99, 97, 49, 52, 101, 57, 102, 98, 56, 53, 48, 102, 100, 56, 101, 98, 101, 53, 100, 99, 101, 50, 48, 34, 44, 34, 116, 114, 97, 99, 107, 95, 110, 117, 109, 98, 101, 114, 34, 58, 34, 49, 54, 98, 102, 102, 56, 51, 97, 101, 55, 100, 49, 52, 101, 54, 57, 98, 53, 57, 101, 51, 54, 101, 54, 57, 53, 54, 51, 51, 99, 48, 49, 34, 44, 34, 101, 110, 116, 114, 121, 34, 58, 34, 109, 97, 120, 105, 109, 101, 34, 44, 34, 100, 101, 108, 105, 118, 101, 114, 121, 34, 58, 123, 34, 110, 97, 109, 101, 34, 58, 34, 89, 101, 115, 115, 101, 110, 105, 97, 32, 66, 111, 101, 104, 109, 34, 44, 34, 112, 104, 111, 110, 101, 34, 58, 34, 56, 52, 51, 45, 50, 49, 57, 45, 49, 48, 55, 54, 34, 44, 34, 122, 105, 112, 34, 58, 34, 51, 54, 49, 49, 49, 34, 44, 34, 99, 105, 116, 121, 34, 58, 34, 77, 111, 110, 116, 103, 111, 109, 101, 114, 121, 34, 44, 34, 97, 100, 100, 114, 101, 115, 115, 34, 58, 34, 50, 48, 50, 52, 32, 77, 101, 114, 114, 105, 108, 121, 32, 68, 114, 105, 118, 101, 34, 44, 34, 114, 101, 103, 105, 111, 110, 34, 58, 34, 65, 76, 34, 44, 34, 101, 109, 97, 105, 108, 34, 58, 34, 75, 90, 117, 110, 76, 86, 83, 64, 117, 116, 115, 113, 83, 109, 106, 46, 110, 101, 116, 34, 125, 44, 34, 112, 97, 121, 109, 101, 110, 116, 34, 58, 123, 34, 116, 114, 97, 110, 115, 97, 99, 116, 105, 111, 110, 34, 58, 34, 97, 54, 48, 51, 56, 55, 55, 98, 54, 48, 100, 100, 52, 100, 56, 51, 57, 50, 51, 48, 100, 53, 102, 57, 55, 99, 54, 48, 48, 102, 56, 57, 34, 44, 34, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 34, 58, 34, 51, 101, 98, 102, 100, 99, 98, 51, 51, 56, 50, 50, 52, 100, 100, 50, 57, 98, 100, 49, 51, 52, 51, 54, 50, 100, 56, 97, 57, 54, 100, 56, 34, 44, 34, 99, 117, 114, 114, 101, 110, 99, 121, 34, 58, 34, 71, 69, 76, 34, 44, 34, 112, 114, 111, 118, 105, 100, 101, 114, 34, 58, 34, 111, 109, 110, 105, 115, 34, 44, 34, 97, 109, 111, 117, 110, 116, 34, 58, 53, 57, 48, 53, 44, 34, 112, 97, 121, 109, 101, 110, 116, 95, 100, 116, 34, 58, 51, 54, 48, 53, 54, 55, 56, 51, 53, 49, 44, 34, 98, 97, 110, 107, 34, 58, 34, 76, 111, 114, 100, 32, 84, 97, 121, 108, 111, 114, 32, 83, 104, 105, 101, 108, 100, 115, 34, 44, 34, 100, 101, 108, 105, 118, 101, 114, 121, 95, 99, 111, 115, 116, 34, 58, 56, 53, 51, 51, 44, 34, 103, 111, 111, 100, 115, 95, 116, 111, 116, 97, 108, 34, 58, 56, 50, 56, 44, 34, 99, 117, 115, 116, 111, 109, 95, 102, 101, 101, 34, 58, 57, 125, 44, 34, 105, 116, 101, 109, 115, 34, 58, 91, 123, 34, 99, 104, 114, 116, 95, 105, 100, 34, 58, 49, 48, 48, 50, 52, 57, 49, 44, 34, 116, 114, 97, 99, 107, 95, 110, 117, 109, 98, 101, 114, 34, 58, 34, 49, 54, 98, 102, 102, 56, 51, 97, 101, 55, 100, 49, 52, 101, 54, 57, 98, 53, 57, 101, 51, 54, 101, 54, 57, 53, 54, 51, 51, 99, 48, 49, 34, 44, 34, 112, 114, 105, 99, 101, 34, 58, 57, 57, 52, 52, 44, 34, 114, 105, 100, 34, 58, 34, 98, 97, 48, 101, 51, 53, 99, 52, 100, 56, 98, 56, 52, 54, 99, 101, 97, 56, 50, 49, 56, 50, 51, 48, 99, 51, 55, 102, 98, 56, 51, 98, 34, 44, 34, 110, 97, 109, 101, 34, 58, 34, 68, 114, 46, 32, 69, 119, 97, 108, 100, 32, 80, 97, 114, 105, 115, 105, 97, 110, 34, 44, 34, 115, 97, 108, 101, 34, 58, 57, 50, 44, 34, 115, 105, 122, 101, 34, 58, 34, 56, 34, 44, 34, 116, 111, 116, 97, 108, 95, 112, 114, 105, 99, 101, 34, 58, 54, 51, 52, 50, 44, 34, 110, 109, 95, 105, 100, 34, 58, 49, 48, 49, 53, 50, 54, 48, 44, 34, 98, 114, 97, 110, 100, 34, 58, 34, 77, 101, 114, 108, 101, 32, 67, 111, 110, 114, 111, 121, 34, 44, 34, 115, 116, 97, 116, 117, 115, 34, 58, 49, 53, 54, 125, 93, 44, 34, 108, 111, 99, 97, 108, 101, 34, 58, 34, 116, 101, 110, 101, 116, 117, 114, 34, 44, 34, 105, 110, 116, 101, 114, 110, 97, 108, 95, 115, 105, 103, 110, 97, 116, 117, 114, 101, 34, 58, 34, 101, 116, 34, 44, 34, 99, 117, 115, 116, 111, 109, 101, 114, 95, 105, 100, 34, 58, 34, 48, 99, 57, 53, 50, 55, 57, 54, 102, 48, 48, 52, 52, 55, 55, 53, 56, 97, 51, 53, 50, 100, 55, 55, 57, 50, 53, 49, 54, 100, 49, 52, 34, 44, 34, 100, 101, 108, 105, 118, 101, 114, 121, 95, 115, 101, 114, 118, 105, 99, 101, 34, 58, 34, 117, 116, 34, 44, 34, 115, 104, 97, 114, 100, 107, 101, 121, 34, 58, 34, 49, 34, 44, 34, 115, 109, 95, 105, 100, 34, 58, 52, 51, 48, 53, 57, 49, 49, 44, 34, 100, 97, 116, 101, 95, 99, 114, 101, 97, 116, 101, 100, 34, 58, 34, 49, 57, 57, 49, 45, 49, 50, 45, 50, 57, 84, 49, 55, 58, 49, 53, 58, 49, 53, 90, 34, 44, 34, 111, 111, 102, 95, 115, 104, 97, 114, 100, 34, 58, 34, 55, 34, 125}
		// dataBytes, _ := json.Marshal(x)
		dataBytes, _ := json.Marshal(data)
		fmt.Printf("data: %+v\n", data)
		fmt.Printf("bytes: %v\n", dataBytes)
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
			CustomFee:    0, //rand.Intn(100),
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
