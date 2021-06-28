package service

import (
	"context"

	"github.com/olvrng/rbot/be/com/flowdef"
	flowdeftypes "github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/com/flowexec"
	"github.com/olvrng/rbot/be/com/flowexec/flowcore"
	"github.com/olvrng/rbot/be/com/flowexec/store"
	"github.com/olvrng/rbot/be/com/flowexec/types"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var _ flowexec.MessengerService = (*MessengerService)(nil)

type MessengerService struct {
	FlowQuery  flowdef.QueryService
	StateStore *store.FlowStateStore
	ActionExec *ActionExecutor
}

func NewMessengerService(
	query flowdef.QueryService,
	stateStore *store.FlowStateStore,
	actionExec *ActionExecutor,
) *MessengerService {
	s := &MessengerService{
		FlowQuery:  query,
		StateStore: stateStore,
		ActionExec: actionExec,
	}
	return s
}

func (s *MessengerService) ReceivedMessage(ctx context.Context, req *types.ReceivedMessageRequest) (*types.ReceivedMessageResponse, error) {
	flowReq := &flowdeftypes.GetFlowByParamRequest{FBPageID: req.PageID}
	flowResp, err := s.FlowQuery.GetFlowByParam(ctx, flowReq)
	if err != nil {
		return nil, xerrors.Errorf(xerrors.NotFound, err, "page_id not found")
	}

	flow, pageID, psid := flowResp.Flow, req.PageID, req.PSID
	state, _, err := loadOrCreateState(ctx, s.StateStore, pageID, psid, flow.ID)
	if err != nil {
		return nil, err
	}

	ex := flowcore.NewExecutor(flow, state)
	stateData := map[string]string{
		"message": req.Message,
	}
	nextState, nextNodes, err := ex.NextState(flowdeftypes.NodeReceivedMessage, stateData)
	if err != nil {
		return nil, err
	}

	actionState := &ActionState{
		PageID: pageID,
		PSID:   psid,
		Extra:  stateData,
	}
	s.ActionExec.ExecuteActions(nextNodes, actionState)

	err = s.StateStore.SaveState(nextState)
	resp := &types.ReceivedMessageResponse{}
	return resp, err
}

func (s *MessengerService) ReceivedPostback(ctx context.Context, req *types.ReceivedPostbackRequest) (*types.ReceivedPostbackResponse, error) {
	flowReq := &flowdeftypes.GetFlowByParamRequest{FBPageID: req.PageID}
	flowResp, err := s.FlowQuery.GetFlowByParam(ctx, flowReq)
	if err != nil {
		return nil, xerrors.Errorf(xerrors.NotFound, err, "page_id not found")
	}

	flow, pageID, psid := flowResp.Flow, req.PageID, req.PSID
	state, _, err := loadOrCreateState(ctx, s.StateStore, pageID, psid, flow.ID)
	if err != nil {
		return nil, err
	}

	ex := flowcore.NewExecutor(flow, state)
	stateData := map[string]string{
		"reply_title":   req.PostbackTitle,
		"reply_payload": req.PostbackPayload,
	}
	nextState, nextNodes, err := ex.NextState(flowdeftypes.NodeReceivedReply, stateData)
	if err != nil {
		return nil, err
	}

	actionState := &ActionState{
		PageID: pageID,
		PSID:   psid,
		Extra:  stateData,
	}
	s.ActionExec.ExecuteActions(nextNodes, actionState)

	err = s.StateStore.SaveState(nextState)
	resp := &types.ReceivedPostbackResponse{}
	return resp, err
}
