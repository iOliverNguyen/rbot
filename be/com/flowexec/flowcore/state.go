package flowcore

import (
	"encoding/json"

	"github.com/olvrng/rbot/be/pkg/dot"
)

type FlowState struct {
	PageID dot.IntID `json:"page_id"`

	PSID dot.IntID `json:"psid"`

	FlowID dot.IntID `json:"flow_id"`

	LastNodeID dot.IntID `json:"last_node_id"`

	NodeID dot.IntID `json:"node_id"`

	Extra map[string]string `json:"extra"`
}

func NewFlowState(pageID, psID, flowID dot.IntID) *FlowState {
	return &FlowState{
		PageID: pageID,
		PSID:   psID,
		FlowID: flowID,
		Extra:  map[string]string{},
	}
}

func (s *FlowState) String() string {
	out, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(out)
}
