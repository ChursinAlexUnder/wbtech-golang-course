package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrderByUid(pool *pgxpool.Pool, ctx context.Context, order_uid string) (Answer, error) {
	var (
		answer                                                    Answer
		order                                                     Orders
		delivery                                                  Delivery
		payment                                                   Payment
		items                                                     []Items
		products                                                  []Product
		err                                                       error
		rowOrder, rowDelivery, rowPayment, rowsItems, rowsProduct pgx.Rows
		orderUidUUID                                              uuid.UUID
	)

	// Парсим для проверки валидности uuid
	orderUidUUID, err = uuid.Parse(order_uid)
	if err != nil {
		return Answer{}, err
	}

	// Получение данных order
	rowOrder, err = pool.Query(ctx, `SELECT *
								FROM orders
								WHERE orders.order_uid = $1`, orderUidUUID)
	if err != nil {
		return Answer{}, err
	}
	defer rowOrder.Close()

	order, err = pgx.CollectOneRow(rowOrder, pgx.RowToStructByName[Orders])
	if err != nil {
		return Answer{}, fmt.Errorf("ошибка форматирования полученной строки order с select запроса в структуру: %v", err)
	}

	// Получение данных delivery
	rowDelivery, err = pool.Query(ctx, `SELECT *
								FROM delivery
								WHERE delivery.uid = $1`, order.Delivery_uid)
	if err != nil {
		return Answer{}, err
	}
	defer rowDelivery.Close()

	delivery, err = pgx.CollectOneRow(rowDelivery, pgx.RowToStructByName[Delivery])
	if err != nil {
		return Answer{}, fmt.Errorf("ошибка форматирования полученной строки delivery с select запроса в структуру: %v", err)
	}

	// Получение данных payment
	rowPayment, err = pool.Query(ctx, `SELECT *
								FROM payment
								WHERE payment.transaction = $1`, order.Payment_uid)
	if err != nil {
		return Answer{}, err
	}
	defer rowPayment.Close()

	payment, err = pgx.CollectOneRow(rowPayment, pgx.RowToStructByName[Payment])
	if err != nil {
		return Answer{}, fmt.Errorf("ошибка форматирования полученной строки payment с select запроса в структуру: %v", err)
	}

	// Получение данных об items
	rowsItems, err = pool.Query(ctx, `SELECT *
								FROM items
								WHERE items.rid = $1`, order.Items_rid)
	if err != nil {
		return Answer{}, err
	}
	defer rowsItems.Close()

	items, err = pgx.CollectRows(rowsItems, pgx.RowToStructByName[Items])
	if err != nil {
		return Answer{}, fmt.Errorf("ошибка форматирования полученных строк items с select запроса в структуру: %v", err)
	}

	// Получение данных об products
	rowsProduct, err = pool.Query(ctx, `SELECT product.*
								FROM product
								JOIN items ON product.nm_id = items.product_id
								WHERE items.rid = $1`, order.Items_rid)
	if err != nil {
		return Answer{}, err
	}
	defer rowsProduct.Close()

	products, err = pgx.CollectRows(rowsProduct, pgx.RowToStructByNameLax[Product])
	if err != nil {
		return Answer{}, fmt.Errorf("ошибка форматирования полученных строк products с select запроса в структуру: %v", err)
	}

	answer = Answer{order, delivery, payment, items, products}

	return answer, nil
}
