package models

type Cart struct {
	ID    int64      `json:"id"`
	Items []CartItem `json:"items"`
}
