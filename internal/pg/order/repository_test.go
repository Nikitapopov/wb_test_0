package pg_order_repo

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
	"wb_test_1/internal/logger/mocks"
	"wb_test_1/internal/models/order"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_pg_repository_GetById(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       *order.Order
		wantErr    bool
	}{
		{
			name: "fail get order",
			args: args{
				id: "b563feb7b2b84b6test",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
						where order_uid = $1`,
					)).
					WithArgs("b563feb7b2b84b6test").
					WillReturnError(errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "success get order",
			args: args{
				id: "b563feb7b2b84b6test",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
						where order_uid = $1`,
					)).
					WithArgs("b563feb7b2b84b6test").
					WillReturnRows(sqlmock.
						NewRows([]string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
						AddRow("b563feb7b2b84b6test", "WBILMTESTTRACK", "WBIL",
							[]byte(`{"name": "Test Testov","phone": "+9720000000","zip": "2639809","city": "Kiryat Mozkin","address": "Ploshad Mira 15","region": "Kraiot","email": "test@gmail.com"}`),
							[]byte(`{"transaction": "b563feb7b2b84b6test","request_id": "","currency": "USD","provider": "wbpay","amount": 1817,"payment_dt": 1637907727,"bank": "alpha","delivery_cost": 1500,"goods_total": 317,"custom_fee": 0}`),
							[]byte(`[{"chrt_id": 9934930,"track_number": "WBILMTESTTRACK","price": 453,"rid": "ab4219087a764ae0btest","name": "Mascaras","sale": 30,"size": "0","total_price": 317,"nm_id": 2389212,"brand": "Vivienne Sabo","status": 202}]`),
							"en", "", "test", "meest", "9", 99, "2021-11-26T06:22:19Z", "1",
						),
					)
			},
			want: &order.Order{
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
		},
		{
			name: "fail unmarshal order with invalid json fields type",
			args: args{
				id: "b563feb7b2b84b6test",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
						where order_uid = $1`,
					)).
					WithArgs("b563feb7b2b84b6test").
					WillReturnRows(sqlmock.
						NewRows([]string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
						AddRow("b563feb7b2b84b6test", "WBILMTESTTRACK", "WBIL",
							[]byte(`{"name": 123,"phone": "+9720000000","zip": "2639809","city": "Kiryat Mozkin","address": "Ploshad Mira 15","region": "Kraiot","email": "test@gmail.com"}`),
							[]byte(`{"transaction": "b563feb7b2b84b6test","request_id": "","currency": "USD","provider": "wbpay","amount": 1817,"payment_dt": 1637907727,"bank": "alpha","delivery_cost": 1500,"goods_total": 317}`),
							[]byte(`[{"chrt_id": 9934930,"track_number": "WBILMTESTTRACK","price": 453,"rid": "ab4219087a764ae0btest","name": "Mascaras","sale": 30,"size": "0","total_price": 317,"nm_id": 2389212,"brand": "Vivienne Sabo","status": 202}]`),
							"en", "", "test", "meest", "9", 99, "2021-11-26T06:22:19Z", "1",
						),
					)
			},
			wantErr: true,
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			u := &repository{
				conn:   mockDB,
				logger: logger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("order repo GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order repo GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pg_repository_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(sqlmock.Sqlmock)
		want       []*order.Order
		wantErr    bool
	}{
		{
			name: "success get not empty order list",
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
					`)).
					WillReturnRows(sqlmock.
						NewRows([]string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
						AddRow("b563feb7b2b84b6test", "WBILMTESTTRACK", "WBIL",
							[]byte(`{"name": "Test Testov","phone": "+9720000000","zip": "2639809","city": "Kiryat Mozkin","address": "Ploshad Mira 15","region": "Kraiot","email": "test@gmail.com"}`),
							[]byte(`{"transaction": "b563feb7b2b84b6test","request_id": "","currency": "USD","provider": "wbpay","amount": 1817,"payment_dt": 1637907727,"bank": "alpha","delivery_cost": 1500,"goods_total": 317, "custom_fee": 0}`),
							[]byte(`[{"chrt_id": 9934930,"track_number": "WBILMTESTTRACK","price": 453,"rid": "ab4219087a764ae0btest","name": "Mascaras","sale": 30,"size": "0","total_price": 317,"nm_id": 2389212,"brand": "Vivienne Sabo","status": 202}]`),
							"en", "", "test", "meest", "9", 99, "2021-11-26T06:22:19Z", "1",
						).
						AddRow("ffa3849b33bc425a8401634f0a5cad30", "b0946b88c91c4371bae2161619e821b2", "quaerat",
							[]byte(`{"name": "Burdette Hintz","phone": "752-963-8110","zip": "20001","city": "Washington","address": "1157 1st Street Northwest","region": "DC","email": "IAMHhWu@UUNOvsZ.com"}`),
							[]byte(`{"transaction": "961d3d14f5cd4513ada1cd5fe8d73111","request_id": "7bdcc57659c54683ad56a7a5679bc325","currency": "EGP","provider": "mollitia","amount": 3306,"payment_dt": 6994328811,"bank": "Princess Hannah Deckow","delivery_cost": 2899,"goods_total": 8784, "custom_fee": 1}`),
							[]byte(`[{"chrt_id": 3692436,"track_number": "b0946b88c91c4371bae2161619e821b2","price": 8463,"rid": "666450c2950e427694c455ae92bd6a10","name": "Mrs. Clemmie Nienow","sale": 68,"size": "6","total_price": 3423,"nm_id": 7092469,"brand": "Rhett Adams","status": 257}, {"chrt_id": 9643507,"track_number": "a30780d1dfdf4b19b29088815778c382","price": 8193,"rid": "a56bf2abe8b4451d88e4ddd795b5518f","name": "Mr. Keagan Rutherford","sale": 87,"size": "4","total_price": 5506,"nm_id": 6686290,"brand": "Jeramy Kuvalis","status": 297}]`),
							"voluptatum", "", "b7c4cbe0bfc240d388d0672906dad992", "rerum", "0", 7961728, "2016-01-27T12:05:10Z", "2",
						),
					)
			},
			want: []*order.Order{
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
			},
		},
		{
			name: "success get empty order list",
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
					`)).
					WillReturnRows(sqlmock.
						NewRows([]string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}),
					)
			},
			want: []*order.Order{},
		},
		{
			name: "fail unmarshal order list with invalid json fields type",
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard
						from public.order
					`)).
					WillReturnRows(sqlmock.
						NewRows([]string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
						AddRow("b563feb7b2b84b6test", "WBILMTESTTRACK", "WBIL",
							[]byte(`{"name": 123,"phone": "+9720000000","zip": "2639809","city": "Kiryat Mozkin","address": "Ploshad Mira 15","region": "Kraiot","email": "test@gmail.com"}`),
							[]byte(`{"transaction": "b563feb7b2b84b6test","request_id": "","currency": "USD","provider": "wbpay","amount": 1817,"payment_dt": 1637907727,"bank": "alpha","delivery_cost": 1500,"goods_total": 317, "custom_fee": 0}`),
							[]byte(`[{"chrt_id": 9934930,"track_number": "WBILMTESTTRACK","price": 453,"rid": "ab4219087a764ae0btest","name": "Mascaras","sale": 30,"size": "0","total_price": 317,"nm_id": 2389212,"brand": "Vivienne Sabo","status": 202}]`),
							"en", "", "test", "meest", "9", 99, "2021-11-26T06:22:19Z", "1",
						).
						AddRow("ffa3849b33bc425a8401634f0a5cad30", "b0946b88c91c4371bae2161619e821b2", "quaerat",
							[]byte(`{"name": "Burdette Hintz","phone": "752-963-8110","zip": "20001","city": "Washington","address": "1157 1st Street Northwest","region": "DC","email": "IAMHhWu@UUNOvsZ.com"}`),
							[]byte(`{"transaction": "961d3d14f5cd4513ada1cd5fe8d73111","request_id": "7bdcc57659c54683ad56a7a5679bc325","currency": "EGP","provider": "mollitia","amount": 3306,"payment_dt": 6994328811,"bank": "Princess Hannah Deckow","delivery_cost": 2899,"goods_total": 8784, "custom_fee": 1}`),
							[]byte(`[{"chrt_id": 3692436,"track_number": "b0946b88c91c4371bae2161619e821b2","price": 8463,"rid": "666450c2950e427694c455ae92bd6a10","name": "Mrs. Clemmie Nienow","sale": 68,"size": "6","total_price": 3423,"nm_id": 7092469,"brand": "Rhett Adams","status": 257}, {"chrt_id": 9643507,"track_number": "a30780d1dfdf4b19b29088815778c382","price": 8193,"rid": "a56bf2abe8b4451d88e4ddd795b5518f","name": "Mr. Keagan Rutherford","sale": 87,"size": "4","total_price": 5506,"nm_id": 6686290,"brand": "Jeramy Kuvalis","status": 297}]`),
							"voluptatum", "", "b7c4cbe0bfc240d388d0672906dad992", "rerum", "0", 7961728, "2016-01-27T12:05:10Z", "2",
						),
					)
			},
			wantErr: true,
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			u := &repository{
				conn:   mockDB,
				logger: logger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("order repo GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order repo GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pg_repository_Insert(t *testing.T) {
	type args struct {
		order order.Order
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success insert order",
			args: args{
				order: order.Order{
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
						{
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
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				deliveryArg := order.OrderDelivery{
					Zip:     "36109",
					City:    "Montgomery",
					Name:    "Edythe Kiehn",
					Email:   "NVFjdsU@fFSDbeI.biz",
					Phone:   "310-296-7815",
					Region:  "AL",
					Address: "74 Ranch Drive",
				}
				paymentArg := order.OrderPayment{
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
				}
				itemsArg := order.OrderItems{
					{
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
				}
				mockSQL.
					ExpectExec(regexp.QuoteMeta(`
						insert into public.order (order_uid, track_number, entry, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment, items)
						values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
					`)).
					WithArgs("1399d76651294101aaa00474655a4172", "116c55084287446cbbade79247be55d4", "quis", "est", "et",
						"a7f24b26edc7440b976f617ed1ead1f9", "ut", "1", 4305911, "1991-12-29T17:15:15Z", "7", deliveryArg, paymentArg, itemsArg).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "fail insert order with repetetive order_uid",
			args: args{
				order: order.Order{
					OrderUid:    "1399d76651294101aaa00474655a4172",
					TrackNumber: "116c55084287446cbbade79247be55d4",
					Entry:       "quis",
					Delivery: order.OrderDelivery{
						Zip:     "36109",
						City:    "Montgomery",
						Name:    "Edythe Kiehn",
						Email:   "NVFjdsU@fFSDbeI.biz",
						Phone:   "310-296-7815",
						Region:  "AL",
						Address: "74 Ranch Drive",
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
						{
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
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				deliveryArg := order.OrderDelivery{
					Zip:     "36109",
					City:    "Montgomery",
					Name:    "Edythe Kiehn",
					Email:   "NVFjdsU@fFSDbeI.biz",
					Phone:   "310-296-7815",
					Region:  "AL",
					Address: "74 Ranch Drive",
				}
				paymentArg := order.OrderPayment{
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
				}
				itemsArg := order.OrderItems{
					{
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
				}
				mockSQL.
					ExpectExec(regexp.QuoteMeta(`
						insert into public.order (order_uid, track_number, entry, locale, internal_signature, customer_id,
							delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment, items)
						values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
					`)).
					WithArgs("1399d76651294101aaa00474655a4172", "116c55084287446cbbade79247be55d4", "quis", "est", "et",
						"a7f24b26edc7440b976f617ed1ead1f9", "ut", "1", 4305911, "1991-12-29T17:15:15Z", "7", deliveryArg, paymentArg, itemsArg).
					WillReturnError(errors.New("Duplicate key order_uid"))
			},
			wantErr: true,
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			u := &repository{
				conn:   mockDB,
				logger: logger,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err := u.Insert(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("order repo Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
