package keybaser

import (
	"context"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

// HandlerFunc .
type HandlerFunc func(ctx context.Context, message *kbchat.SubscriptionMessage)

// registerMiddleware .
func (k *Keybaser) registerMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, message *kbchat.SubscriptionMessage) {
		k.responseLog(next)(ctx, message)
	}
}

// responseLog .
func (k *Keybaser) responseLog(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, message *kbchat.SubscriptionMessage) {
		startTime := time.Now()

		next(ctx, message)

		k.logger.
			WithField("severity", "info").
			WithField("duration", time.Since(startTime).String()).
			WithField("username", message.Message.Sender.Username).
			WithField("channel", message.Message.Channel.Name).
			WithField("topic", message.Message.Channel.TopicName).
			Infoln("received a message")
	}
}
