package service

import (
	"context"

	"github.com/olvrng/rbot/be/com/flowdef"
	"github.com/olvrng/rbot/be/com/flowdef/store"
	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var _ flowdef.QueryService = (*FlowQueryService)(nil)

type FlowQueryService struct {
	Store store.FlowStore
}

func NewFlowQueryService(flowStore store.FlowStore) *FlowQueryService {
	s := &FlowQueryService{
		Store: flowStore,
	}
	return s
}

func (f *FlowQueryService) GetFlowByID(ctx context.Context, req *types.GetFlowByIDRequest) (*types.FlowResponse, error) {
	flow := f.Store.LoadFlowByID(req.ID)
	if flow == nil {
		return nil, xerrors.Errorf(xerrors.NotFound, nil, "not found")
	}
	return &types.FlowResponse{Flow: flow}, nil
}

func (f *FlowQueryService) GetFlowByParam(ctx context.Context, req *types.GetFlowByParamRequest) (*types.FlowResponse, error) {

	if req.FBPageID == 0 {
		return nil, xerrors.Errorf(xerrors.NotFound, nil, "not found")
	}

	flow := f.Store.LoadFlowByPageID(req.FBPageID)
	if flow == nil {
		return nil, xerrors.Errorf(xerrors.NotFound, nil, "not found")
	}
	return &types.FlowResponse{Flow: flow}, nil
}
