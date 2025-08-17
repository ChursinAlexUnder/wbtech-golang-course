package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SelectOrderByUid(ctx context.Context, pool *pgxpool.Pool, order_uid string) (Orders, error) {
	var (
		order                                        Orders
		err                                          error
		rowOrder, rowDelivery, rowPayment, rowsItems pgx.Rows
		orderUidUUID                                 uuid.UUID
	)

	// Парсим для проверки валидности uuid
	orderUidUUID, err = uuid.Parse(order_uid)
	if err != nil {
		return Orders{}, err
	}

	// Получение данных order
	rowOrder, err = pool.Query(ctx, `SELECT *
								FROM orders
								WHERE orders.order_uid = $1`, orderUidUUID)
	if err != nil {
		return Orders{}, err
	}
	defer rowOrder.Close()

	if rowOrder.Next() {
		if err := rowOrder.Scan(
			&order.Order_uid,
			&order.Track_number,
			&order.Entry,
			&order.Delivery_uid,
			&order.Payment_uid,
			&order.Locale,
			&order.Internal_signature,
			&order.Customer_id,
			&order.Delivery_service,
			&order.Shardkey,
			&order.Sm_id,
			&order.Date_created,
			&order.Oof_shard); err != nil {
			return Orders{}, err
		}
	} else {
		return Orders{}, fmt.Errorf("order with this order_uid was not found")
	}

	// Получение данных delivery
	rowDelivery, err = pool.Query(ctx, `SELECT *
								FROM delivery
								WHERE delivery.uid = $1`, order.Delivery_uid)
	if err != nil {
		return Orders{}, err
	}
	defer rowDelivery.Close()

	order.Delivery, err = pgx.CollectOneRow(rowDelivery, pgx.RowToStructByName[Delivery])
	if err != nil {
		return Orders{}, fmt.Errorf("ошибка взятия данных с таблицы delivery с базы данных: %v", err)
	}

	// Получение данных payment
	rowPayment, err = pool.Query(ctx, `SELECT *
								FROM payment
								WHERE payment.transaction = $1`, order.Payment_uid)
	if err != nil {
		return Orders{}, err
	}
	defer rowPayment.Close()

	order.Payment, err = pgx.CollectOneRow(rowPayment, pgx.RowToStructByName[Payment])
	if err != nil {
		return Orders{}, fmt.Errorf("ошибка взятия данных с таблицы payment с базы данных: %v", err)
	}

	// Получение данных об items
	rowsItems, err = pool.Query(ctx, `SELECT *
								FROM items
								WHERE items.order_uid = $1`, order.Order_uid)
	if err != nil {
		return Orders{}, err
	}
	defer rowsItems.Close()

	order.Items, err = pgx.CollectRows(rowsItems, pgx.RowToStructByName[Items])
	if err != nil {
		return Orders{}, fmt.Errorf("ошибка взятия данных с таблицы items с базы данных: %v", err)
	}

	return order, nil
}

// Получение самых свежих 1000 заказов для кеша
func SelectOrdersForCache(ctx context.Context, pool *pgxpool.Pool) ([]Orders, error) {
	var (
		orders                                       []Orders = make([]Orders, 100, 1000)
		order                                        Orders
		err                                          error
		rowOrder, rowDelivery, rowPayment, rowsItems pgx.Rows
	)
	// Получение актуальных данных order
	rowOrder, err = pool.Query(ctx, `SELECT *
								FROM orders
								ORDER BY date_created DESC
								LIMIT 1000`)
	if err != nil {
		return []Orders{}, err
	}
	defer rowOrder.Close()

	for rowOrder.Next() {
		if err := rowOrder.Scan(
			&order.Order_uid,
			&order.Track_number,
			&order.Entry,
			&order.Delivery_uid,
			&order.Payment_uid,
			&order.Locale,
			&order.Internal_signature,
			&order.Customer_id,
			&order.Delivery_service,
			&order.Shardkey,
			&order.Sm_id,
			&order.Date_created,
			&order.Oof_shard); err != nil {
			return []Orders{}, err
		}

		// Взятие информации текущего заказа из остальных таблиц
		// Получение данных delivery
		rowDelivery, err = pool.Query(ctx, `SELECT *
								FROM delivery
								WHERE delivery.uid = $1`, order.Delivery_uid)
		if err != nil {
			return []Orders{}, err
		}

		order.Delivery, err = pgx.CollectOneRow(rowDelivery, pgx.RowToStructByName[Delivery])
		if err != nil {
			return []Orders{}, fmt.Errorf("ошибка взятия данных с таблицы delivery с базы данных: %v", err)
		}
		rowDelivery.Close()

		// Получение данных payment
		rowPayment, err = pool.Query(ctx, `SELECT *
								FROM payment
								WHERE payment.transaction = $1`, order.Payment_uid)
		if err != nil {
			return []Orders{}, err
		}

		order.Payment, err = pgx.CollectOneRow(rowPayment, pgx.RowToStructByName[Payment])
		if err != nil {
			return []Orders{}, fmt.Errorf("ошибка взятия данных с таблицы payment с базы данных: %v", err)
		}
		rowPayment.Close()

		// Получение данных об items
		rowsItems, err = pool.Query(ctx, `SELECT *
								FROM items
								WHERE items.order_uid = $1`, order.Order_uid)
		if err != nil {
			return []Orders{}, err
		}

		order.Items, err = pgx.CollectRows(rowsItems, pgx.RowToStructByName[Items])
		if err != nil {
			return []Orders{}, fmt.Errorf("ошибка взятия данных с таблицы items с базы данных: %v", err)
		}
		rowsItems.Close()

		// Добавляем собранный заказ в итоговый срез
		orders = append(orders, order)
	}
	return orders, nil
}
