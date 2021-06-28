package types

import "github.com/olvrng/rbot/be/com/integration/fbmsg"

type ReceivedMessageRequest struct {
	PageID  fbmsg.IntID `json:"page_id"`
	PSID    fbmsg.IntID `json:"psid"`
	Message string      `json:"message"`
}

type ReceivedMessageResponse struct {
}

type ReceivedPostbackRequest struct {
	PageID          fbmsg.IntID `json:"page_id"`
	PSID            fbmsg.IntID `json:"psid"`
	PostbackTitle   string      `json:"postback_title"`
	PostbackPayload string      `json:"postback_payload"`
}

type ReceivedPostbackResponse struct {
}
