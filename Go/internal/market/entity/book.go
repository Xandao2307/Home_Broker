package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order         []*Order
	Trasactions   []*Trasaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(ordersChan chan *Order, ordersChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Trasactions:   []*Trasaction{},
		OrdersChan:    ordersChan,
		OrdersChanOut: ordersChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChan {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)
			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)
				if sellOrder.PedingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- sellOrder
					b.OrdersChanOut <- order
					if sellOrder.PedingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)
			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				buyOrder := buyOrders.Pop().(*Order)
				if buyOrder.PedingShares > 0 {
					transaction := NewTransaction(buyOrder, order, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- order
					if buyOrder.PedingShares > 0 {
						buyOrders.Push(buyOrders)
					}
				}
			}
		}

	}
}

func (b *Book) AddTransaction(transaction *Trasaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PedingShares
	buyingShares := transaction.BuyingOrder.PedingShares

	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, -minShares)

	transaction.AddBuyOrderPedingShares(-minShares)
	transaction.AddSellOrderPedingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.ClosedBuyOrder()
	transaction.ClosedSellOrder()

	b.Trasactions = append(b.Trasactions, transaction)
}
