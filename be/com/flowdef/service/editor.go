package service

import (
	"context"

	"github.com/olvrng/rbot/be/com/flowdef"
	"github.com/olvrng/rbot/be/com/flowdef/store"
	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var _ flowdef.EditorService = (*FlowEditorService)(nil)

type FlowEditorService struct {
	Store store.FlowStore
}

func NewFlowEditorService(flowStore store.FlowStore) *FlowEditorService {
	s := &FlowEditorService{
		Store: flowStore,
	}
	return s
}

func (s *FlowEditorService) CreateFlow(ctx context.Context, req *types.CreateFlowRequest) (*types.CreateFlowResponse, error) {
	flow := s.Store.SaveFlow(req.Flow)
	return &types.CreateFlowResponse{Flow: flow}, nil
}

func (s *FlowEditorService) UpdateFlow(ctx context.Context, req *types.CreateFlowRequest) (*types.CreateFlowResponse, error) {
	if req.Flow.ID == 0 {
		return nil, xerrors.Errorf(xerrors.InvalidArgument, nil, "id is required")
	}
	flow := s.Store.SaveFlow(req.Flow)
	return &types.CreateFlowResponse{Flow: flow}, nil
}
