package fbmsg

// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events

type WebhookMessage struct {
	Object string          `json:"object"`
	Entry  []*WebhookEntry `json:"entry"`
}

type WebhookEntry struct {
	ID        IntID             `json:"id"`
	Time      Timestamp         `json:"time"`
	Messaging []*EntryMessaging `json:"messaging"`
}

type EntryMessaging struct {
	Sender    SenderID      `json:"sender"`
	Recipient RecipientID   `json:"recipient"`
	Timestamp Timestamp     `json:"timestamp"`
	Message   *MessageData  `json:"message,omitempty"`
	Postback  *PostbackData `json:"postback,omitempty"`
	Reaction  *ReactionData `json:"reaction,omitempty"`
	Read      *ReadData     `json:"read,omitempty"`
	Delivery  *DeliveryData `json:"delivery,omitempty"`
}

type SenderID struct {
	ID      IntID  `json:"id,omitempty"`
	UserRef string `json:"user_ref,omitempty"`
}

type RecipientID struct {
	ID IntID `json:"id"`
}

type MessageData struct {
	MID         StrID                `json:"mid,omitempty"`
	AppID       IntID                `json:"app_id,omitempty"`
	Text        string               `json:"text,omitempty"`
	IsEcho      bool                 `json:"is_echo,omitempty"`
	QuickReply  *QuickReplyData      `json:"quick_reply,omitempty"`
	ReplyTo     *ReplyToData         `json:"reply_to,omitempty"`
	Attachments []*AttachmentItem    `json:"attachments,omitempty"`
	Referral    *MessageReferralData `json:"referral,omitempty"`
}

type PostbackData struct {
	MID      StrID                 `json:"mid,omitempty"`
	Title    string                `json:"title,omitempty"`
	Payload  string                `json:"payload,omitempty"`
	Referral *PostbackReferralData `json:"referral,omitempty"`
}

type ReactionData struct {
	MID      StrID  `json:"mid,omitempty"`
	Action   string `json:"action,omitempty"`
	Emoji    string `json:"emoji,omitempty"`
	Reaction string `json:"reaction,omitempty"`
}

type ReadData struct {
	Watermark IntID `json:"watermark,omitempty"`
}

type DeliveryData struct {
	MIDs []StrID `json:"mids"`
}

type QuickReplyData struct {
	Payload string `json:"payload,omitempty"`
}

type ReplyToData struct {
	MID StrID `json:"mid,omitempty"`
}

type AttachmentItem struct {
	Type    AttachmentType         `json:"type,omitempty"`
	Payload *AttachmentPayloadData `json:"payload,omitempty"`
}

type AttachmentPayloadData struct {
	URL string `json:"url,omitempty"`
}

type MessageReferralData struct {
}

type PostbackReferralData struct {
	Ref         string `json:"ref"`
	Source      string `json:"source"`
	Type        string `json:"type"`
	RefererURI  string `json:"referer_uri"`
	IsGuestUser string `json:"is_guest_user"`
}
