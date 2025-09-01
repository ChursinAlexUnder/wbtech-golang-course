package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertOrder(ctx context.Context, pool *pgxpool.Pool, order Orders) error {
	var err error

	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// Вставляем в delivery
	_, err = pool.Exec(ctx,
		`INSERT INTO delivery (uid, name, phone, zip, city, email, address, region) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.Delivery.Uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Email, order.Delivery.Address, order.Delivery.Region)
	if err != nil {
		return err
	}

	// Вставляем в orders
	_, err = pool.Exec(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, delivery_uid, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		order.Order_uid, order.Track_number, order.Entry, order.Delivery_uid, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		return err
	}

	// Вставляем в payment
	_, err = pool.Exec(ctx,
		`INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.Payment.Transaction, order.Payment.Request_id, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt, order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total, order.Payment.Custom_fee)
	if err != nil {
		return err
	}

	// Вставляем записи в items
	for _, item := range order.Items {
		_, err = pool.Exec(ctx,
			`INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_price, item.Nm_id, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	// Коммитим
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
