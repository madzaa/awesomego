package models

type UserBalance struct {
	ID      uint64 `json:"userId"`
	Balance string `json:"balance"`
}
