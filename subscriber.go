package keybaser

import "github.com/keybase/go-keybase-chat-bot/kbchat"

// NewSubscription .
type NewSubscription interface {
	Read() (kbchat.SubscriptionMessage, error)
}

// Subscriber .
type Subscriber interface {
	Subscribe() (NewSubscription, error)
}

type subscriber struct {
	client KeybaseChatAPIClient
}

func newSubscriber(client KeybaseChatAPIClient) Subscriber {
	return &subscriber{client: client}
}

// Subscribe event
func (s *subscriber) Subscribe() (NewSubscription, error) {
	return s.client.ListenForNewTextMessages()
}
