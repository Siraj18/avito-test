package models

import "time"

type Transaction struct {
	Id        string    `json:"id" db:"id"`
	ToId      *string   `json:"to_id,omitempty" db:"to_id"`
	FromId    *string   `json:"from_id,omitempty" db:"from_id"`
	Money     float64   `json:"money" db:"money"`
	Operation string    `json:"operation" db:"operation"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TransactionGetQuery struct {
	Id string `json:"id"`
}

type AllTransactionsGetQuery struct {
	Id       string `json:"id"`
	SortType string `json:"sort_type"`
	Limit    int    `json:"limit"`
	Page     int    `json:"page"`
}
