package transaction

type TransactionRequest struct {
	EmailReceiver string `json:"email_receiver"`
	AmountTransaction int `json:"amount_transaction"`
}