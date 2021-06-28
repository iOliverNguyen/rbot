package types

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/olvrng/rbot/be/pkg/dot"
)

type NodeType string

const (
	NodeCompletedOrder  = "trigger:completed_order"
	NodeReceivedMessage = "trigger:received_message"
	NodeReceivedReply   = "trigger:received_reply"
	NodeSendMessage     = "action:send_message"
)

type Flow struct {
	ID      dot.IntID   `json:"id"`
	PageIDs []dot.IntID `json:"page_ids"`
	Nodes   []*Node     `json:"nodes"`
}

func (f *Flow) NodeByID(nodeID dot.IntID) *Node {
	for _, node := range f.Nodes {
		if node.ID == nodeID {
			return node
		}
	}
	return nil
}

type Node struct {
	ID      dot.IntID    `json:"id,omitempty"`
	Payload *NodePayload `json:"payload,omitempty"`
}

func (n *Node) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type NodePayload struct {
	// Type            NodeType `json:"type"`
	CompletedOrder  *CompletedOrderNodeData
	SendMessage     *SendMessageNodeData
	ReceivedMessage *ReceivedMessageNodeData
}

func (n *NodePayload) Type() NodeType {
	if n.CompletedOrder != nil {
		return NodeCompletedOrder
	}
	if n.SendMessage != nil {
		return NodeSendMessage
	}
	if n.ReceivedMessage != nil {
		return NodeReceivedMessage
	}
	return ""
}

func (n *NodePayload) MarshalJSON() ([]byte, error) {
	switch {
	case n.CompletedOrder != nil:
		n.CompletedOrder.Type = NodeCompletedOrder
		return json.Marshal(n.CompletedOrder)
	case n.SendMessage != nil:
		n.SendMessage.Type = NodeSendMessage
		return json.Marshal(n.SendMessage)
	case n.ReceivedMessage != nil:
		n.ReceivedMessage.Type = NodeReceivedMessage
		return json.Marshal(n.ReceivedMessage)
	default:
		return nil, nil
	}
}

func (n *NodePayload) UnmarshalJSON(data []byte) error {
	var t struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	switch t.Type {
	case NodeCompletedOrder:
		return json.Unmarshal(data, &n.CompletedOrder)
	case NodeSendMessage:
		return json.Unmarshal(data, &n.SendMessage)
	case NodeReceivedMessage:
		return json.Unmarshal(data, &n.ReceivedMessage)
	default:
		return errors.New("unknown node")
	}
}

func (n *NodePayload) Next(typ NodeType, data map[string]string) dot.IntID {
	switch {
	case n.CompletedOrder != nil:
		return n.CompletedOrder.Next(typ, data)
	case n.SendMessage != nil:
		return n.SendMessage.Next(typ, data)
	case n.ReceivedMessage != nil:
		return n.ReceivedMessage.Next(typ, data)
	default:
		return 0
	}
}

type NodeNextInterface interface {
	Next(typ NodeType, data map[string]string) dot.IntID
}

type CompletedOrderNodeData struct {
	Type   NodeType  `json:"type"`
	NextID dot.IntID `json:"next_id,omitempty"`
	Fields []string  `json:"fields"`
}

func (n *CompletedOrderNodeData) Next(typ NodeType, data map[string]string) dot.IntID {
	return n.NextID
}

type SendMessageNodeData struct {
	Type         NodeType          `json:"type"`
	Template     string            `json:"template"`
	NextID       dot.IntID         `json:"next_id,omitempty"`
	QuickReplies []*QuickReplyItem `json:"quick_replies"`
}

func (n *SendMessageNodeData) Next(typ NodeType, data map[string]string) dot.IntID {
	switch typ {
	case NodeReceivedMessage:
		return n.NextID

	case NodeReceivedReply:
		for _, reply := range n.QuickReplies {
			if reply.Code == data["reply_payload"] {
				return reply.NextID
			}
		}
		return n.NextID

	default:
		return 0
	}
}

type QuickReplyItem struct {
	Text   string    `json:"text"`
	Code   string    `json:"code"`
	NextID dot.IntID `json:"next_id"`
}

type ReceivedMessageNodeData struct {
	Type   NodeType  `json:"type"`
	NextID dot.IntID `json:"next_id"`
}

func (n *ReceivedMessageNodeData) Next(typ NodeType, data map[string]string) dot.IntID {
	switch typ {
	case NodeReceivedMessage:
		return n.NextID
	}
	return 0
}
