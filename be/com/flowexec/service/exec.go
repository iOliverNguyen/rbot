package service

import (
	"context"
	"sync"
	"time"

	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/com/integration/fbmsg"
	"github.com/olvrng/rbot/be/pkg/dot"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var ll = l.New()
var ls = ll.Sugar()

type ActionState struct {
	PageID dot.IntID
	PSID   dot.IntID
	Extra  map[string]string
}

type ActionExecutor struct {
	FBClient *fbmsg.Client
}

func NewActionExecutor(fbClient *fbmsg.Client) *ActionExecutor {
	ex := &ActionExecutor{FBClient: fbClient}
	return ex
}

func (ex *ActionExecutor) ExecuteActions(nodes []*types.Node, state *ActionState) {
	ls.Debug("execute actions: ", nodes)
	if len(nodes) == 0 {
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func(node *types.Node) {
			defer wg.Done()
			_ = ex.ExecuteAction(ctx, node, state)
		}(node)
	}
	wg.Wait()
	ctxCancel()
}

func (ex *ActionExecutor) ExecuteAction(ctx context.Context, node *types.Node, state *ActionState) (_err error) {
	defer func() {
		if re := recover(); re != nil {
			_err = xerrors.Errorf(xerrors.Internal, nil, "recovered")
		}
		if _err != nil {
			ll.Error("execute action", l.Error(_err))
		}
	}()

	switch node.Payload.Type() {
	case types.NodeSendMessage:
		return ex.execSendMessage(ctx, node, state)

	default:
		ls.Error("unknown node type ", node)
		return xerrors.Errorf(xerrors.Internal, nil, "unknown node type")
	}
}

func (ex *ActionExecutor) execSendMessage(ctx context.Context, node *types.Node, state *ActionState) error {
	payload := node.Payload.SendMessage

	respMsg := &fbmsg.SendMessageData{}
	if len(payload.QuickReplies) == 0 {
		respMsg.Text = payload.Template

	} else {
		buttons := make([]*fbmsg.ButtonItem, 0, len(payload.QuickReplies))
		for _, reply := range payload.QuickReplies {
			btn := fbmsg.PostbackButtonData{
				Type:    fbmsg.ButtonTypePostback,
				Title:   reply.Text,
				Payload: reply.Code,
			}.Wrap()
			buttons = append(buttons, btn)
		}
		respMsg.Attachment = &fbmsg.SendAttachmentData{
			Type: fbmsg.SendAttachmentTypeTemplate,
			Payload: &fbmsg.PayloadData{
				TemplateType: fbmsg.TemplateTypeGeneric,
				Elements: []*fbmsg.ElementItem{
					{
						Title:         payload.Template,
						Subtitle:      "",
						DefaultAction: nil,
						Buttons:       buttons,
					},
				},
			},
		}
	}

	sendReq := &fbmsg.SendRequest{
		Recipient: &fbmsg.SendRecipientData{ID: state.PSID},
		Message:   respMsg,
	}
	return ex.FBClient.CallSendAPI(ctx, sendReq)
}
