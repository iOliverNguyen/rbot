package store

import (
	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/pkg/dot"
)

type FlowStore interface {
	SaveFlow(flow *types.Flow) *types.Flow

	LoadFlowByID(id dot.IntID) *types.Flow

	LoadFlowByPageID(pageID dot.IntID) *types.Flow
}
