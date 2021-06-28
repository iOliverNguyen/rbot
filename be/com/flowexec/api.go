package flowexec

import (
	"context"

	"github.com/olvrng/rbot/be/com/flowexec/types"
)

// +gen:api

// +api:path=/api/flow/exec/messenger
type MessengerService interface {
	ReceivedMessage(ctx context.Context, req *types.ReceivedMessageRequest) (*types.ReceivedMessageResponse, error)

	ReceivedPostback(ctx context.Context, req *types.ReceivedPostbackRequest) (*types.ReceivedPostbackResponse, error)
}

// +api:path=/api/flow/exec/order
type OrderService interface {
	ReceivedCompletedOrder(ctx context.Context, req *types.ReceivedCompletedOrderRequest) (*types.ReceivedCompletedOrderResponse, error)
}
