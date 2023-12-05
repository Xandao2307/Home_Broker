package entity

type Order struct {
	ID           string
	Investor     *Investor
	Asset        *Asset
	Shares       int
	PedingShares int
	Price        float64
	OrderType    string
	Status       string
	Transactions []*Trasaction
}

func NewOrder(orderId string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		ID:           orderId,
		Investor:     investor,
		Asset:        asset,
		Shares:       shares,
		PedingShares: shares,
		Price:        price,
		OrderType:    orderType,
		Status:       "Open",
		Transactions: []*Trasaction{},
	}
}
