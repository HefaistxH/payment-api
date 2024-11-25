package dto

type Payment struct {
	CustomerId string  `json:"customer_id"`
	MerchantId string  `json:"merchant_id"`
	Pin        string  `json:"pin"`
	Amount     float64 `json:"amount"`
	Message    string  `json:"message"`
}
