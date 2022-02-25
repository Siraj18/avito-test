package models

type User struct {
	Id      string  `json:"id" db:"id"`
	Balance float32 `json:"balance" db:"balance"`
}

type UserGetBalanceQuery struct {
	Id string `json:"id"`
}
