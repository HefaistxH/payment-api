package entity

import "time"

type History struct {
	Id         string    `json:"id"`
	CustomerId string    `json:"user_id"`
	MerchantId string    `json:"merchant_id"`
	Activity   string    `json:"activity"`
	Amount     float64   `json:"amount"`
	Message    string    `json:"message"`
	TimeStamp  time.Time `json:"time_stamp"`
}
