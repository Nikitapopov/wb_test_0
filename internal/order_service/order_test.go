package order_service

import (
	"errors"
	"reflect"
	"testing"
	cacheMocks "wb_test_1/internal/cache/order/mocks"
	loggerMocks "wb_test_1/internal/logger/mocks"
	"wb_test_1/internal/models/order"
	pgMocks "wb_test_1/internal/pg/order/mocks"

	"github.com/golang/mock/gomock"
)

func Test_order_SyncCacheWithDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orders := []*order.Order{
		{
			OrderUid:    "b563feb7b2b84b6test",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: order.OrderDelivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: order.OrderPayment{
				Transaction:  "b563feb7b2b84b6test",
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: order.OrderItems{
				{
					ChrtId:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmId:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmId:              99,
			DateCreated:       "2021-11-26T06:22:19Z",
			OofShard:          "1",
		},
		{
			OrderUid:    "ffa3849b33bc425a8401634f0a5cad30",
			TrackNumber: "b0946b88c91c4371bae2161619e821b2",
			Entry:       "quaerat",
			Delivery: order.OrderDelivery{
				Name:    "Burdette Hintz",
				Phone:   "752-963-8110",
				Zip:     "20001",
				City:    "Washington",
				Address: "1157 1st Street Northwest",
				Region:  "DC",
				Email:   "IAMHhWu@UUNOvsZ.com",
			},
			Payment: order.OrderPayment{
				Transaction:  "961d3d14f5cd4513ada1cd5fe8d73111",
				RequestId:    "7bdcc57659c54683ad56a7a5679bc325",
				Currency:     "EGP",
				Provider:     "mollitia",
				Amount:       3306,
				PaymentDt:    6994328811,
				Bank:         "Princess Hannah Deckow",
				DeliveryCost: 2899,
				GoodsTotal:   8784,
				CustomFee:    1,
			},
			Items: order.OrderItems{
				{
					ChrtId:      3692436,
					TrackNumber: "b0946b88c91c4371bae2161619e821b2",
					Price:       8463,
					Rid:         "666450c2950e427694c455ae92bd6a10",
					Name:        "Mrs. Clemmie Nienow",
					Sale:        68,
					Size:        "6",
					TotalPrice:  3423,
					NmId:        7092469,
					Brand:       "Rhett Adams",
					Status:      257,
				},
				{
					ChrtId:      9643507,
					TrackNumber: "a30780d1dfdf4b19b29088815778c382",
					Price:       8193,
					Rid:         "a56bf2abe8b4451d88e4ddd795b5518f",
					Name:        "Mr. Keagan Rutherford",
					Sale:        87,
					Size:        "4",
					TotalPrice:  5506,
					NmId:        6686290,
					Brand:       "Jeramy Kuvalis",
					Status:      297,
				},
			},
			Locale:            "voluptatum",
			InternalSignature: "",
			CustomerId:        "b7c4cbe0bfc240d388d0672906dad992",
			DeliveryService:   "rerum",
			Shardkey:          "0",
			SmId:              7961728,
			DateCreated:       "2016-01-27T12:05:10Z",
			OofShard:          "2",
		},
	}

	type args struct {
		input order.Order
	}
	tests := []struct {
		name       string
		beforeTest func(repo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker)
		want       error
	}{
		{
			name: "success sync all orders",
			beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				gomock.InOrder(
					dbRepo.EXPECT().
						GetAll().
						Return(orders, nil).
						Times(1),
					cacheRepo.EXPECT().
						Insert(orders[0].OrderUid, *orders[0]).Return(nil).
						Times(1),
					cacheRepo.EXPECT().
						Insert(orders[1].OrderUid, *orders[1]).Return(nil).
						Times(1),
				)
			},
			want: nil,
		},
		{
			name: "success sync empty list of orders",
			beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				gomock.InOrder(
					dbRepo.EXPECT().
						GetAll().
						Return([]*order.Order{}, nil).
						Times(1),
					cacheRepo.EXPECT().
						Insert(orders[0].OrderUid, *orders[0]).Return(nil).
						Times(0),
					cacheRepo.EXPECT().
						Insert(orders[1].OrderUid, *orders[1]).Return(nil).
						Times(0),
				)
			},
			want: nil,
		},
		{
			name: "fail sync orders, because dbRepo.GetAll() error",
			beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				dbRepo.EXPECT().
					GetAll().
					Return(nil, errors.New("Incorrect query to db")).
					Times(1)
			},
			want: errors.New("Incorrect query to db"),
		},
		{
			name: "fail sync orders, because cacheRepo.Insert() error",
			beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				gomock.InOrder(
					dbRepo.EXPECT().
						GetAll().
						Return(orders, nil).
						Times(1),
					cacheRepo.EXPECT().
						Insert(orders[0].OrderUid, *orders[0]).Return(errors.New("some cache error")).
						Times(1),
					cacheRepo.EXPECT().
						Insert(orders[1].OrderUid, *orders[1]).Return(nil).
						Times(0),
				)
			},
			want: errors.New("some cache error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPgRepo := pgMocks.NewMockPgDbRepositoryWorker(ctrl)
			mockCacheRepo := cacheMocks.NewMockCacheRepositoryWorker(ctrl)
			mockLogger := loggerMocks.NewMockLogger()

			w := &orderService{
				mockPgRepo,
				mockCacheRepo,
				mockLogger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockPgRepo, mockCacheRepo)
			}

			err := w.syncCacheWithDb()
			if (err != nil) != (tt.want != nil) {
				t.Errorf("order_service.syncCacheWithDb() error = %v, wantErr %v", err, tt.want)
				return
			}
		})

	}
}

func Test_order_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testOrder := &order.Order{
		OrderUid:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: order.OrderDelivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: order.OrderPayment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: order.OrderItems{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}

	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(cacheRepo *cacheMocks.MockCacheRepositoryWorker)
		want       *order.Order
		wantErr    error
	}{
		{
			name: "success get order by id",
			beforeTest: func(cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				cacheRepo.EXPECT().
					GetById(testOrder.OrderUid).Return(testOrder, nil).
					Times(1)
			},
			want: testOrder,
		},
		{
			name: "fail get order by id, because it isn't found",
			beforeTest: func(cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				cacheRepo.EXPECT().
					GetById(testOrder.OrderUid).Return(nil, errors.New("order not found")).
					Times(1)
			},
			wantErr: errors.New("order not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPgRepo := pgMocks.NewMockPgDbRepositoryWorker(ctrl)
			mockCacheRepo := cacheMocks.NewMockCacheRepositoryWorker(ctrl)
			mockLogger := loggerMocks.NewMockLogger()

			service := &orderService{
				mockPgRepo,
				mockCacheRepo,
				mockLogger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockCacheRepo)
			}

			got, err := service.GetById(testOrder.OrderUid)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("order_service.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order repo GetById() = %v, want %v", got, tt.want)
			}
		})

	}
}

func Test_order_GetIdsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ordersIds := []string{"0173deae63b7417ca9ab93f514da5735", "02198275cf934ee3a494808f250f994e", "02cbfa10ed9a473182ef56e561abf652"}

	tests := []struct {
		name       string
		beforeTest func(cacheRepo *cacheMocks.MockCacheRepositoryWorker)
		want       []string
		wantErr    error
	}{
		{
			name: "success get not empty orders ids list",
			beforeTest: func(cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				cacheRepo.EXPECT().
					GetIdsList().Return(ordersIds).
					Times(1)
			},
			want: ordersIds,
		},
		{
			name: "success get empty orders ids list",
			beforeTest: func(cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				cacheRepo.EXPECT().
					GetIdsList().Return([]string{}).
					Times(1)
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPgRepo := pgMocks.NewMockPgDbRepositoryWorker(ctrl)
			mockCacheRepo := cacheMocks.NewMockCacheRepositoryWorker(ctrl)
			mockLogger := loggerMocks.NewMockLogger()

			service := &orderService{
				mockPgRepo,
				mockCacheRepo,
				mockLogger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockCacheRepo)
			}

			got := service.GetIdsList()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order repo GetById() = %v, want %v", got, tt.want)
			}
		})

	}
}

func Test_order_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		input order.Order
	}

	testOrder := order.Order{
		OrderUid:    "1399d76651294101aaa00474655a4172",
		TrackNumber: "116c55084287446cbbade79247be55d4",
		Entry:       "quis",
		Delivery: order.OrderDelivery{
			Name:    "Edythe Kiehn",
			Phone:   "310-296-7815",
			Zip:     "36109",
			City:    "Montgomery",
			Address: "74 Ranch Drive",
			Region:  "AL",
			Email:   "NVFjdsU@fFSDbeI.biz",
		},
		Payment: order.OrderPayment{
			Bank:         "Prof. Viva Carroll",
			Amount:       296,
			Currency:     "XAG",
			Provider:     "ad",
			CustomFee:    3,
			PaymentDt:    4864056167,
			RequestId:    "4b883c803756413cafb97d9db93324e1",
			GoodsTotal:   1019,
			Transaction:  "ad0a3973d3ad482abf160e1e14af5ff8",
			DeliveryCost: 4811,
		},
		Items: order.OrderItems{
			order.OrderItem{
				Rid:         "28d358188f904342b880c87c4bd69ef2",
				Name:        "Miss Alanis Witting",
				Sale:        59,
				Size:        "2",
				Brand:       "Ruby Gislason",
				NmId:        1934617,
				Price:       3115,
				Status:      188,
				ChrtId:      4570773,
				TotalPrice:  2954,
				TrackNumber: "116c55084287446cbbade79247be55d4",
			},
		},
		Locale:            "est",
		InternalSignature: "et",
		CustomerId:        "a7f24b26edc7440b976f617ed1ead1f9",
		DeliveryService:   "ut",
		Shardkey:          "1",
		SmId:              4305911,
		DateCreated:       "1991-12-29T17:15:15Z",
		OofShard:          "7",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(repo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker)
		want       error
	}{
		// {
		// 	name: "success inserting new order",
		// 	args: args{
		// 		input: testOrder,
		// 	},
		// 	beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
		// 		gomock.InOrder(
		// 			dbRepo.EXPECT().Insert(testOrder).Return(nil).Times(1),
		// 			cacheRepo.EXPECT().Insert(testOrder.OrderUid, testOrder).Return(nil).Times(1),
		// 		)
		// 	},
		// 	want: nil,
		// },
		{
			name: "fail inserting existing order",
			args: args{
				input: testOrder,
			},
			beforeTest: func(dbRepo *pgMocks.MockPgDbRepositoryWorker, cacheRepo *cacheMocks.MockCacheRepositoryWorker) {
				gomock.InOrder(
					dbRepo.EXPECT().
						Insert(testOrder).
						Return(errors.New("Already exists")).
						Times(1),
					cacheRepo.EXPECT().
						Insert(testOrder.OrderUid, testOrder).
						Return(errors.New("Already exists")).
						Times(0),
				)
			},
			want: errors.New("Already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPgRepo := pgMocks.NewMockPgDbRepositoryWorker(ctrl)
			mockCacheRepo := cacheMocks.NewMockCacheRepositoryWorker(ctrl)
			mockLogger := loggerMocks.NewMockLogger()

			service := &orderService{
				mockPgRepo,
				mockCacheRepo,
				mockLogger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockPgRepo, mockCacheRepo)
			}

			err := service.Insert(tt.args.input)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("order_service.Insert() error = %v, wantErr %v", err, tt.want)
				return
			}
		})

	}
}
