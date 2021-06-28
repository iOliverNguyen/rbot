package types

import "github.com/olvrng/rbot/be/pkg/dot"

type ReceivedCompletedOrderRequest struct {
	PageID dot.IntID `json:"page_id"`

	// mock: psid_orderid
	OrderID string `json:"order_id"`

	Desc string `json:"string"`

	Amount int `json:"amount"`
}

type ReceivedCompletedOrderResponse struct {
}
