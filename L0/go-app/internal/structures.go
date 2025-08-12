package internal

import (
	"time"

	"github.com/google/uuid"
)

type Orders struct {
	Order_uid          uuid.UUID `json:"order_uid"`
	Track_number       string    `json:"track_number"`
	Entry              string    `json:"entry"`
	Delivery           `json:"delivery"`
	Payment            `json:"payment"`
	Items              []Items   `json:"items"`
	Locale             string    `json:"locale"`
	Internal_signature string    `json:"internal_signature"`
	Customer_id        string    `json:"customer_id"`
	Delivery_service   string    `json:"delivery_service"`
	Shardkey           string    `json:"shardkey"`
	Sm_id              int       `json:"sm_id"`
	Date_created       time.Time `json:"date_created"`
	Oof_shard          string    `json:"oof_shard"`
}

type Delivery struct {
	Uid     uuid.UUID `json:"uid"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Zip     string    `json:"zip"`
	City    string    `json:"city"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
	Region  string    `json:"region"`
}

type Payment struct {
	Transaction   uuid.UUID `json:"transaction"`
	Request_id    string    `json:"request_id"`
	Currency      string    `json:"currency"`
	Provider      string    `json:"provider"`
	Amount        float64   `json:"amount"`
	Payment_dt    int       `json:"payment_dt"`
	Bank          string    `json:"bank"`
	Delivery_cost float64   `json:"delivery_cost"`
	Goods_total   int       `json:"goods_total"`
	Custom_fee    float64   `json:"custom_fee"`
}

type Items struct {
	Items_uid    uuid.UUID `json:"items_uid"`
	Track_number string    `json:"track_number"`
	Rid          uuid.UUID `json:"rid"`
	Status       int       `json:"status"`
	Product      `json:"product"`
}

type Product struct {
	Nm_id       int     `json:"nm_id"`
	Chrt_id     int     `json:"chrt_id"`
	Price       float64 `json:"price"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	Total_price float64 `json:"total_price"`
	Brand       string  `json:"brand"`
}
