package keybaser

import (
	"context"
	"fmt"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// KeybaseChatAPIClient .
type KeybaseChatAPIClient interface {
	ListenForNewTextMessages() (kbchat.NewSubscription, error)
	SendMessage(channel chat1.ChatChannel, body string) (kbchat.SendResponse, error)
	GetUsername() string
}

// ResponseConstructor .
type ResponseConstructor func(channel chat1.ChatChannel, client KeybaseChatAPIClient) ResponseWriter

// Keybaser contains the Keybase API
type Keybaser struct {
	client              KeybaseChatAPIClient
	botCommands         []BotCommand
	logger              logrus.FieldLogger
	middlewareFunc      func(next HandlerFunc) HandlerFunc
	responseConstructor ResponseConstructor
	subscriber          Subscriber
}

// ClientOption is option client initialize
type ClientOption func(*Keybaser)

// WithLogger is option for set custom logger
func WithLogger(logger logrus.FieldLogger) ClientOption {
	return func(bot *Keybaser) {
		bot.logger = logger
	}
}

// SetMiddlewareFunc is option for set custom middleware
func SetMiddlewareFunc(middlewareFunc func(next HandlerFunc) HandlerFunc) ClientOption {
	return func(bot *Keybaser) {
		bot.middlewareFunc = middlewareFunc
	}
}

// SetResponseConstructor is option for set custom response constructor
func SetResponseConstructor(constructor ResponseConstructor) ClientOption {
	return func(bot *Keybaser) {
		bot.responseConstructor = constructor
	}
}

// SetSubscriber is option for set custom subscriber
func SetSubscriber(subscriber Subscriber) ClientOption {
	return func(bot *Keybaser) {
		bot.subscriber = subscriber
	}
}

// New creates a new client using the Keybase API
func New(kc KeybaseChatAPIClient, opts ...ClientOption) (*Keybaser, error) {
	bot := &Keybaser{
		client:              kc,
		logger:              newDefaultLogger(),
		responseConstructor: NewResponse,
		subscriber:          newSubscriber(kc),
	}

	bot.middlewareFunc = bot.registerMiddleware

	for _, opt := range opts {
		opt(bot)
	}

	return bot, nil
}

// Client returns the internal keybase chat API
func (k *Keybaser) Client() KeybaseChatAPIClient {
	return k.client
}

// Command define a new command and append it to the list of existing commands
func (k *Keybaser) Command(usage string, definition *CommandDefinition) {
	k.botCommands = append(k.botCommands, NewBotCommand(usage, definition))
}

func (k *Keybaser) handleMessage(ctx context.Context, message *kbchat.SubscriptionMessage) {
	resp := k.responseConstructor(message.Message.Channel, k.client)

	for _, cmd := range k.botCommands {
		parameters, isMatch := cmd.Match(message.Message.Content.Text.Body)
		if !isMatch {
			continue
		}

		req := NewRequest(ctx, message, parameters)
		if cmd.Definition().AuthorizationFunc != nil && !cmd.Definition().AuthorizationFunc(req) {
			resp.ReportError(errors.New("You are not authorized to execute this command"))
			k.logger.Infoln("authorization error")
			return
		}

		cmd.Execute(req, resp)
	}
}

func (k *Keybaser) handler() HandlerFunc {
	return k.middlewareFunc(k.handleMessage)
}

func (k *Keybaser) handle(ctx context.Context, message *kbchat.SubscriptionMessage) {
	k.handler()(ctx, message)
}

// Listen receives events from Keybase and each is handled as needed
func (k *Keybaser) Listen(ctx context.Context) error {
	k.Command("help", k.helpCommand())

	k.logger.WithField("severity", "info").Infof("starting subscribe messages as %s", k.client.GetUsername())

	sub, err := k.subscriber.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := sub.Read()
			if err != nil {
				k.logger.Errorf("failed to read message: %s", err)
				continue
			}

			if k.isFromBot(&msg) {
				continue
			}

			// in channel, mention is required
			// in case of direct message, mention is NOT required
			if !k.isBotMentioned(&msg) && !k.isDirectMessage(&msg) {
				continue
			}

			go k.handle(ctx, &msg)
		}
	}
}

func (k *Keybaser) isFromBot(message *kbchat.SubscriptionMessage) bool {
	return message.Message.Sender.Username == k.client.GetUsername()
}

func (k *Keybaser) isBotMentioned(message *kbchat.SubscriptionMessage) bool {
	return strings.Contains(message.Message.Content.Text.Body, "@"+k.client.GetUsername())
}

func (k *Keybaser) isDirectMessage(message *kbchat.SubscriptionMessage) bool {
	participants := strings.Split(message.Message.Channel.Name, ",")
	return len(participants) == 2
}
