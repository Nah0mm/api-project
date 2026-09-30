package model

type Job struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}
