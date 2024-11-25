package entity

type Customer struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
