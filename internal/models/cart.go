package models

type Cart struct {
	ID    string     `json:"id"`
	Items []CartItem `json:"items"`
}
