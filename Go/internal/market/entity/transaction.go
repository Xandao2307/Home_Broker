package entity

import (
	"time"

	"github.com/google/uuid"
)

type Trasaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Trasaction {
	total := float64(shares) * price

	return &Trasaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (t *Trasaction) CalculateTotal(shares int, price float64) {
	t.Total = float64(shares) * price
}

func (t *Trasaction) AddBuyOrderPedingShares(minShares int) {
	t.BuyingOrder.PedingShares += minShares
}

func (t *Trasaction) AddSellOrderPedingShares(minShares int) {
	t.SellingOrder.PedingShares += minShares
}

func (t *Trasaction) ClosedBuyOrder() {
	if t.BuyingOrder.PedingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}

func (t *Trasaction) ClosedSellOrder() {
	if t.SellingOrder.PedingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}
