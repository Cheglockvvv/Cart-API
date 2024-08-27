package repository

import "github.com/jmoiron/sqlx"

type CartItem struct {
	DB *sqlx.DB
}

func InitCartItem(db *sqlx.DB) *CartItem {
	cartItem := &CartItem{DB: db}

	return cartItem
}
