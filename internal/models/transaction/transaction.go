package transaction

type TransactionRequest struct {
	Amount int `json:"amount"`
	To string `json:"to"`
}