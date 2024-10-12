package handler

type PaymentResponse struct {
	Amount  int    `json:"amount"`
	SnapURL string `json:"snap_url"`
}

type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}
