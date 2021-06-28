package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/olvrng/rbot/be/com/flowdef"
	flowdeftypes "github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/com/flowexec"
	"github.com/olvrng/rbot/be/com/flowexec/flowcore"
	"github.com/olvrng/rbot/be/com/flowexec/store"
	"github.com/olvrng/rbot/be/com/flowexec/types"
	"github.com/olvrng/rbot/be/pkg/dot"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var _ flowexec.OrderService = (*OrderService)(nil)

type OrderService struct {
	FlowQuery  flowdef.QueryService
	StateStore *store.FlowStateStore
	ActionExec *ActionExecutor
}

func NewOrderService(
	query flowdef.QueryService,
	stateStore *store.FlowStateStore,
	actionExec *ActionExecutor,
) *OrderService {
	s := &OrderService{
		FlowQuery:  query,
		StateStore: stateStore,
		ActionExec: actionExec,
	}
	return s
}

func (s *OrderService) ReceivedCompletedOrder(ctx context.Context, req *types.ReceivedCompletedOrderRequest) (*types.ReceivedCompletedOrderResponse, error) {
	flowReq := &flowdeftypes.GetFlowByParamRequest{FBPageID: req.PageID}
	flowResp, err := s.FlowQuery.GetFlowByParam(ctx, flowReq)
	if err != nil {
		ll.Error("page_id not found", l.ID("page_id", req.PageID))
		return nil, xerrors.Errorf(xerrors.NotFound, err, "page_id not found")
	}
	flow := flowResp.Flow

	psid, err := mockGetPSIDFromOrderID(req.OrderID)
	if err != nil {
		return nil, xerrors.Errorf(xerrors.NotFound, err, "psid not found")
	}
	state, _, err := loadOrCreateState(ctx, s.StateStore, req.PageID, psid, flow.ID)
	if err != nil {
		return nil, err
	}

	ex := flowcore.NewExecutor(flow, state)
	stateData := map[string]string{
		"order_id": req.OrderID,
		"desc":     req.Desc,
		"amount":   fmt.Sprint(req.Amount),
	}
	nextState, nextNodes, err := ex.NextState(flowdeftypes.NodeCompletedOrder, stateData)
	if err != nil {
		return nil, err
	}

	actionState := &ActionState{
		PageID: req.PageID,
		PSID:   psid,
		Extra:  stateData,
	}
	s.ActionExec.ExecuteActions(nextNodes, actionState)

	err = s.StateStore.SaveState(nextState)
	resp := &types.ReceivedCompletedOrderResponse{}
	return resp, err
}

func mockGetPSIDFromOrderID(orderID string) (dot.IntID, error) {
	parts := strings.Split(orderID, "_")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	return dot.IntID(id), err
}

func loadOrCreateState(
	ctx context.Context,
	stateStore *store.FlowStateStore,
	pageID, psid, flowID dot.IntID,
) (state *flowcore.FlowState, isNew bool, _ error) {
	state, err := stateStore.LoadState(ctx, pageID, psid)
	switch xerrors.GetCode(err) {
	case xerrors.NoError:
		return state, false, nil

	case xerrors.NotFound:
		state = flowcore.NewFlowState(pageID, psid, flowID)
		return state, true, nil

	default:
		return nil, false, xerrors.Errorf(xerrors.Internal, err, "internal error")
	}
}
