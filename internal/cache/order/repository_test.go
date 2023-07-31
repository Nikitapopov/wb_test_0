package order_repo_cache

import (
	"fmt"
	"reflect"
	"testing"
	"wb_test_1/internal/logger/mocks"
	"wb_test_1/internal/models/order"
	"wb_test_1/pkg/go_cache"
)

func Test_cache_repository_GetById(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		want    *order.Order
		wantErr bool
	}{
		{
			name: "fail get order",
			args: args{
				id: "nonexistent_item_id",
			},
			wantErr: true,
		},
		{
			name: "success get order",
			args: args{
				id: "existing_item_id",
			},
			wantErr: true,
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoCacheClient := go_cache.NewMockGoCacher()

			repo := &repository{
				client: mockGoCacheClient,
				logger: logger,
			}

			got, err := repo.GetById(tt.args.id)
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

func Test_cache_repository_GetIdsList(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name: "success get orders ids list",
			want: []string{"id_1", "id_2", "id_3"},
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoCacheClient := go_cache.NewMockGoCacher()

			repo := &repository{
				client: mockGoCacheClient,
				logger: logger,
			}

			got := repo.GetIdsList()
			fmt.Print("got:")
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order repo GetIdsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_repository_Insert(t *testing.T) {
	type args struct {
		order order.Order
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
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
			wantErr: false,
		},
		{
			name: "fail insert order with existed order_uid",
			args: args{
				order: order.Order{
					OrderUid:    "already_existed_id",
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
			wantErr: true,
		},
	}

	logger := mocks.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoCacheClient := go_cache.NewMockGoCacher()

			repo := &repository{
				client: mockGoCacheClient,
				logger: logger,
			}

			err := repo.Insert(tt.args.order.OrderUid, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("order repo GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
