package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/olvrng/rbot/be/com/flowexec/service"
	"github.com/olvrng/rbot/be/com/flowexec/types"
	"github.com/olvrng/rbot/be/com/integration/fbmsg"
	"github.com/olvrng/rbot/be/pkg/l"
)

var ll = l.New()
var ls = ll.Sugar()

type WebhookService struct {
	Client           *fbmsg.Client
	VerifyToken      string
	MessengerService *service.MessengerService
}

func NewWebhookService(
	client *fbmsg.Client, token string,
	messengerService *service.MessengerService,
) *WebhookService {
	s := &WebhookService{
		Client:           client,
		VerifyToken:      token,
		MessengerService: messengerService,
	}
	return s
}

func (s *WebhookService) HandleVerification(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		ll.Error("webhook: bad request (parse form)")
		w.WriteHeader(400)
		return
	}

	mode := req.Form.Get("hub.mode")
	token := req.Form.Get("hub.verify_token")
	challenge := req.Form.Get("hub.challenge")
	if mode == "" || token == "" {
		ll.Error("webhook: bad request")
		w.WriteHeader(400)
		return
	}

	if mode == "subscribe" && token == s.VerifyToken {
		ll.Info("webhook: verify token", l.String("challenge", challenge))
		_, _ = fmt.Fprint(w, challenge)
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(403)
}

// When you receive a webhook event, you must always return a 200 OK HTTP
// response. The Messenger Platform will resend the webhook event every 20
// seconds, until a 200 OK response is received. Failing to return a 200 OK may
// cause your webhook to be unsubscribed by the Messenger Platform.

func (s *WebhookService) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	var r io.Reader = req.Body
	if ce := ll.Check(l.DebugLevel, "debugging"); ce != nil {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ll.Error("can not read request", l.Error(err))
			w.WriteHeader(400)
			return
		}
		r = bytes.NewReader(data)
		ls.Debugf("webhook body: %s", data)
	}

	var body fbmsg.WebhookMessage
	if err := json.NewDecoder(r).Decode(&body); err != nil {
		_, _ = fmt.Fprintf(w, "can not decode json")
		w.WriteHeader(400)
		return
	}

	var ctx = req.Context()
	var err error
	switch body.Object {
	case "page":
		for _, entry := range body.Entry {
			if len(entry.Messaging) == 0 {
				continue
			}
			pageID := entry.ID
			// Get the webhook event. entry.messaging is an array, but
			// will only ever contain one event, so we get index 0
			event := entry.Messaging[0]
			senderPSID := event.Sender.ID
			ll.Debug("received webhook", l.ID("senderPSID", senderPSID))

			switch {
			case event.Message != nil:
				err = s.HandleMessage(ctx, pageID, event.Sender, event.Message)

			case event.Postback != nil:
				err = s.HandlePostback(ctx, pageID, event.Sender, event.Postback)

			default:
				ll.Debug("webhook: ignore message", l.ID("entry.id", entry.ID))
			}
		}

	default:
		ll.Debug("webhook: ignore object", l.String("object", body.Object))
	}

	if err != nil {
		ll.Error("webhook: error", l.Error(err))
	} else {
		ll.Info("webhook: ok")
	}
	// always return a 200 OK HTTP
	w.WriteHeader(200)
}

func (s *WebhookService) HandleMessage(ctx context.Context, pageID fbmsg.IntID, sender fbmsg.SenderID, msg *fbmsg.MessageData) error {
	if msg.IsEcho {
		ll.Debug("webhook: echo message, ignore")
		return nil
	}

	req := &types.ReceivedMessageRequest{
		PageID:  pageID,
		PSID:    sender.ID,
		Message: msg.Text,
	}
	_, err := s.MessengerService.ReceivedMessage(ctx, req)
	return err
}

func (s *WebhookService) HandlePostback(ctx context.Context, pageID fbmsg.IntID, sender fbmsg.SenderID, msg *fbmsg.PostbackData) error {

	req := &types.ReceivedPostbackRequest{
		PageID:          pageID,
		PSID:            sender.ID,
		PostbackTitle:   msg.Title,
		PostbackPayload: msg.Payload,
	}
	_, err := s.MessengerService.ReceivedPostback(ctx, req)
	return err
}
