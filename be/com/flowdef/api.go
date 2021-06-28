package flowdef

import (
	"context"

	"github.com/olvrng/rbot/be/com/flowdef/types"
)

// +gen:api

// +api:path=/api/flow/def/editor
type EditorService interface {
	CreateFlow(ctx context.Context, req *types.CreateFlowRequest) (*types.CreateFlowResponse, error)

	UpdateFlow(ctx context.Context, req *types.CreateFlowRequest) (*types.CreateFlowResponse, error)
}

// +api:path=/api/flow/def/query
type QueryService interface {
	GetFlowByID(ctx context.Context, req *types.GetFlowByIDRequest) (*types.FlowResponse, error)

	GetFlowByParam(ctx context.Context, req *types.GetFlowByParamRequest) (*types.FlowResponse, error)
}
