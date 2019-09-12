package keybaser

import (
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
)

// A ResponseWriter interface is used to respond to an event
type ResponseWriter interface {
	Reply(text string)
	ReportError(err error)
	Client() KeybaseChatAPIClient
}

// NewResponse creates a new response structure
func NewResponse(channel chat1.ChatChannel, client KeybaseChatAPIClient) ResponseWriter {
	return &response{channel: channel, client: client}
}

type response struct {
	channel chat1.ChatChannel
	client  KeybaseChatAPIClient
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *response) ReportError(err error) {
	_, _ = r.client.SendMessage(r.channel, err.Error())
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string) {
	_, _ = r.client.SendMessage(r.channel, message)
}

// Client returns the slack client
func (r *response) Client() KeybaseChatAPIClient {
	return r.client
}
