package pg_order_repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"wb_test_1/internal/logger"
	"wb_test_1/internal/order"
)

const dbTimeout = time.Second * 3

type repository struct {
	conn   *sql.DB
	logger logger.Logger
}

func NewRepo(db *sql.DB, logger logger.Logger) order.DbRepositoryWorker {
	return &repository{
		conn:   db,
		logger: logger,
	}
}

func (r *repository) GetById(id string) (*order.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		from public.order
		where order_uid = $1
	`

	var order order.Order
	row := r.conn.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&order.OrderUid,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery,
		&order.Payment,
		&order.Items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		err := fmt.Errorf("receiving order: %w", err)
		r.logger.Error(err)
		return nil, err
	}

	return &order, nil
}

func (r *repository) GetAll() ([]*order.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		from public.order
	`

	rows, err := r.conn.QueryContext(ctx, stmt)
	if err != nil {
		err := fmt.Errorf("receiving orders: %w", err)
		r.logger.Error(err)
		return nil, err
	}

	orders := make([]*order.Order, 0)
	for rows.Next() {
		var order order.Order
		err := rows.Scan(
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Items,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			err := fmt.Errorf("scanning of order: %w", err)
			r.logger.Error(err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *repository) Insert(order order.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		insert into public.order (order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment, items)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	execArgs := []any{
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
		order.Delivery,
		order.Payment,
		order.Items,
	}

	_, err := r.conn.ExecContext(ctx, stmt, execArgs...)
	if err != nil {
		err := fmt.Errorf("order insert: %w", err)
		r.logger.Error(err)
		return err
	}

	return nil
}
