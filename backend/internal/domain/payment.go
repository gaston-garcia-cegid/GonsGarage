package domain

import "time"

type Payment struct {
	ID        string
	Amount    float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPayment(id string, amount float64, currency string) *Payment {
	return &Payment{
		ID:        id,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
func (p *Payment) UpdateAmount(amount float64) {
	p.Amount = amount
	p.UpdatedAt = time.Now()
}
func (p *Payment) UpdateCurrency(currency string) {
	p.Currency = currency
	p.UpdatedAt = time.Now()
}
