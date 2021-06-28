package flowcore

import (
	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var ll = l.New()
var ls = ll.Sugar()

type Executor struct {
	Flow  *types.Flow
	State *FlowState
}

func NewExecutor(flow *types.Flow, state *FlowState) *Executor {
	ex := &Executor{
		Flow:  flow,
		State: state,
	}
	return ex
}

func (ex *Executor) NextState(nodeType types.NodeType, data map[string]string) (_nextState *FlowState, _nextNodes []*types.Node, _err error) {

	defer func() {
		if ce := ll.Check(l.DebugLevel, "debug next state"); ce != nil {
			ls.Debugf("currState=%v", ex.State)
			ls.Debugf("nextState: err=%v nodeType=%v nextState=%v nextNodes=%v", _err, nodeType, _nextState, _nextNodes)
		}
	}()

	flow, state := ex.Flow, ex.State
	node := flow.NodeByID(state.NodeID)
	if node != nil {
		nextState, nextNodes, ok := ex.execNextNodes(state, node, nodeType, data)
		if ok {
			return nextState, nextNodes, nil
		}
	}

	// fallback when no node detected
	node, fallback := mockGetNodeByType(ex.Flow, nodeType)
	if node == nil {
		node = fallback
	}
	if node != nil {
		nextState, nextNodes, ok := ex.execNextNodes(state, node, nodeType, data)
		if ok {
			return nextState, nextNodes, nil
		}
	}

	return nil, nil, xerrors.Errorf(xerrors.Aborted, nil, "can not execute state")
}

func (ex *Executor) execNextNodes(state *FlowState, node *types.Node, nodeType types.NodeType, data map[string]string) (_nextState *FlowState, _nodes []*types.Node, ok bool) {

	flow := ex.Flow
	nextNodeID := node.Payload.Next(nodeType, data)
	if nextNodeID != 0 {
		nextNode := flow.NodeByID(nextNodeID)
		if nextNode == nil {
			return nil, nil, false
		}

		nextState := NewFlowState(state.PageID, state.PSID, state.FlowID)
		nextState.NodeID = nextNodeID
		nextState.Extra = data

		_nodes = append(_nodes, nextNode)
		return nextState, _nodes, true
	}
	return nil, nil, false
}

func mockGetNodeByType(flow *types.Flow, nodeType types.NodeType) (node, fallback *types.Node) {
	for _, nd := range flow.Nodes {
		// fallback node
		if nd.Payload.Type() == types.NodeReceivedMessage && fallback == nil {
			fallback = nd
		}
		// real node
		if nd.Payload.Type() == nodeType {
			return nd, nil
		}
	}
	return nil, fallback
}
