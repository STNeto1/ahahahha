package core

import "time"

type Category struct {
	ID        string  `db:"id"`
	Name      string  `db:"name"`
	Slug      string  `db:"slug"`
	CreatedAt string  `db:"created_at"`
	ParentID  *string `db:"parent_id"`
}

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
	Quantity    uint32    `db:"quantity"`
	SellerID    string    `db:"seller_id"`
	CategoryID  string    `db:"category_id"`
	CreatedAt   string    `db:"created_at"`
}
