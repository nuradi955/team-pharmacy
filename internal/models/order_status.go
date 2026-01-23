package models

type OrderStatus string

const (
	OrderStatusDraft          OrderStatus = "draft"
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusCanceled       OrderStatus = "canceled"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusCompleted      OrderStatus = "completed"
)

var allowedOrderStatusTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPendingPayment: {OrderStatusPaid, OrderStatusCanceled},
	OrderStatusPaid:           {OrderStatusShipped},
	OrderStatusShipped:        {OrderStatusCompleted},
}

func CanChangeOrderStatus(from, to OrderStatus) bool {
	next, ok := allowedOrderStatusTransitions[from]
	if !ok {
		return false
	}
	for _, s := range next {
		if s == to {
			return true
		}
	}
	return false
}
