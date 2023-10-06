package models

import "time"

type Item struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Level       *int32    `db:"level"`
	Rarity      int       `db:"rarity"`
	Image       *string   `db:"image"`
	FromTime    time.Time `db:"from_time"`
	ToTime      time.Time `db:"to_time"`
	Price       float64   `db:"price"`
	BuyoutPrice float64   `db:"buyout_price"`
	Quantity    int       `db:"quantity"`
	CreatedAt   string    `db:"created_at"`

	CategoryID string `db:"category_id" gqlgen:"-"`
	SellerID   string `db:"seller_id" gqlgen:"-"`
}
