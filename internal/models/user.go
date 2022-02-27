package models

type User struct {
	Id      string  `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}

type UserGetBalanceQuery struct {
	Id string `json:"id"`
}

type UserChangeBalanceQuery struct {
	Id    string  `json:"id"`
	Money float64 `json:"money"`
}

type UserTransferBalanceQuery struct {
	FromId string  `json:"from_id"`
	ToId   string  `json:"to_id"`
	Money  float64 `json:"money"`
}
