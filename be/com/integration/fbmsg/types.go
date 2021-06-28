package fbmsg

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/olvrng/rbot/be/pkg/dot"
)

// https://developers.facebook.com/docs/messenger-platform/reference/send-api

type URL string
type StrID string
type IntID = dot.IntID
type Timestamp = dot.Timestamp

type ObjectType string
type SenderAction string
type ContentType string
type AttachmentType string
type SendAttachmentType string
type TemplateType string
type ButtonType string

const (
	ObjectTypePage ObjectType = "page"

	SenderActionTypingOn  SenderAction = "typing_on"
	SenderActionTypingOff SenderAction = "typing_off"
	SenderActionMarkSeen  SenderAction = "mark_seen"

	ContentTypeText  ContentType = "text"
	ContentTypePhone ContentType = "user_phone_number"
	ContentTypeEmail ContentType = "user_email"

	AttachmentTypeAudio    AttachmentType = "audio"
	AttachmentTypeFile     AttachmentType = "file"
	AttachmentTypeImage    AttachmentType = "image"
	AttachmentTypeLocation AttachmentType = "location"
	AttachmentTypeVideo    AttachmentType = "video"
	AttachmentTypeFallback AttachmentType = "fallback"

	SendAttachmentTypeImage    SendAttachmentType = "image"
	SendAttachmentTypeAudio    SendAttachmentType = "audio"
	SendAttachmentTypeVideo    SendAttachmentType = "video"
	SendAttachmentTypeFile     SendAttachmentType = "file"
	SendAttachmentTypeTemplate SendAttachmentType = "template"

	TemplateTypeGeneric TemplateType = "generic"
	TemplateTypeButton  TemplateType = "button"

	ButtonTypeURL      ButtonType = "web_url"
	ButtonTypePostback ButtonType = "postback"
)

// SendRequest implements Messenger API
type SendRequest struct {
	MessagingType    string             `json:"messaging_type,omitempty"`
	Recipient        *SendRecipientData `json:"recipient,omitempty"`
	Message          *SendMessageData   `json:"message,omitempty"`
	SenderAction     string             `json:"sender_action,omitempty"`
	NotificationType string             `json:"notification_type,omitempty"`
	Tag              string             `json:"tag,omitempty"`
}

type SendRecipientData struct {
	ID        IntID  `json:"id,omitempty"`
	UserRef   string `json:"user_ref,omitempty"`
	PostId    string `json:"post_id,omitempty"`
	CommentID string `json:"comment_id,omitempty"`
}

type SendMessageData struct {
	Text         string              `json:"text,omitempty"`
	Attachment   *SendAttachmentData `json:"attachment,omitempty"`
	QuickReplies []*QuickReplyItem   `json:"quick_replies,omitempty"`
	Metadata     string              `json:"metadata,omitempty"`
}

type SendAttachmentData struct {
	Type    SendAttachmentType `json:"type,omitempty"`
	Payload *PayloadData       `json:"payload,omitempty"`
}

type QuickReplyItem struct {
	ContentType string       `json:"content_type,omitempty"`
	Title       string       `json:"title,omitempty"`
	Payload     *PayloadData `json:"payload,omitempty"`
	ImageURL    string       `json:"image_url,omitempty"`
}

type SendResponse struct {
	RecipientID IntID  `json:"recipient_id,omitempty"`
	MessageID   string `json:"message_id,omitempty"`
}

// PayloadData can be template or file attachment. We only implement 2 templates here.
//
// Template: generic, button, ...
//
// https://developers.facebook.com/docs/messenger-platform/reference/templates#available_templates
// https://developers.facebook.com/docs/messenger-platform/reference/templates/generic
// https://developers.facebook.com/docs/messenger-platform/reference/templates/button
type PayloadData struct {
	TemplateType TemplateType   `json:"template_type"`
	Elements     []*ElementItem `json:"elements"`
}

type ElementItem struct {
	Title         string             `json:"title,omitempty"`
	Subtitle      string             `json:"subtitle,omitempty"`
	ImageURL      string             `json:"image_url,omitempty"`
	DefaultAction *DefaultActionData `json:"default_action,omitempty"`
	Buttons       []*ButtonItem      `json:"buttons,omitempty"`
}

// DefaultActionData accepts the same properties as URL button, except title.
type DefaultActionData URLButtonData

type ButtonItem struct {
	URLButton      *URLButtonData
	PostbackButton *PostbackButtonData
}

func (b ButtonItem) MarshalJSON() ([]byte, error) {
	switch {
	case b.URLButton != nil:
		return json.Marshal(b.URLButton)
	case b.PostbackButton != nil:
		return json.Marshal(b.PostbackButton)
	default:
		return nil, errors.New("no data")
	}
}

func (b *ButtonItem) UnmarshalJSON(data []byte) error {
	var d struct {
		Type ButtonType `json:"type"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	switch d.Type {
	case ButtonTypeURL:
		return json.Unmarshal(data, &b.URLButton)
	case ButtonTypePostback:
		return json.Unmarshal(data, &b.PostbackButton)
	default:
		return nil // ignore
	}
}

// URLButtonData
// https://developers.facebook.com/docs/messenger-platform/reference/buttons/url
type URLButtonData struct {
	Type  ButtonType `json:"type"`
	Title string     `json:"title"`
	URL   string     `json:"url"`
}

func (d URLButtonData) Wrap() *ButtonItem {
	d.Type = ButtonTypeURL
	return &ButtonItem{URLButton: &d}
}

// PostbackButtonData
// https://developers.facebook.com/docs/messenger-platform/reference/buttons/postback
type PostbackButtonData struct {
	Type    ButtonType `json:"type"`
	Title   string     `json:"title"`
	Payload string     `json:"payload"`
}

func (d PostbackButtonData) Wrap() *ButtonItem {
	d.Type = ButtonTypePostback
	return &ButtonItem{PostbackButton: &d}
}
