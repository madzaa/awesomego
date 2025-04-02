package models

type Transaction struct {
	ID         string `json:"transactionId"`
	UserID     uint64 `json:"userId"`
	Amount     uint64 `json:"amount"`
	State      string `json:"state"`
	SourceType string `json:"-"`
}

type TransactionRequest struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionId string `json:"transactionId"`
}
