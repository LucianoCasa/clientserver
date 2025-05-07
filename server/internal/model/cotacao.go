package model

import "time"

type Cotacao struct {
	ID        int
	Valor     float64
	CreatedAt time.Time
}
