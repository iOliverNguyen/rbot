package types

import (
	"github.com/olvrng/rbot/be/pkg/dot"
)

type CreateFlowRequest struct {
	Flow *Flow `json:"flow"`
}

type CreateFlowResponse struct {
	Flow *Flow `json:"flow"`
}

type UpdateFlowRequest struct {
	Flow *Flow `json:"flow"`
}

type UpdateFlowResponse struct {
	Flow *Flow `json:"flow"`
}

type GetFlowByIDRequest struct {
	ID dot.IntID `json:"id"`
}

type GetFlowByParamRequest struct {
	FBPageID dot.IntID `json:"id"`
}

type FlowResponse struct {
	Flow *Flow `json:"flow"`
}
