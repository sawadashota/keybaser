package keybaser_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/sawadashota/keybaser"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type testNewSubscription struct {
	message *kbchat.SubscriptionMessage
}

func newTestNewSubscription(message *kbchat.SubscriptionMessage) *testNewSubscription {
	return &testNewSubscription{
		message: message,
	}
}

func (t *testNewSubscription) Read() (kbchat.SubscriptionMessage, error) {
	return *t.message, nil
}

type testSubscriber struct {
	message *kbchat.SubscriptionMessage
}

func newTestSubscriber(message *kbchat.SubscriptionMessage) *testSubscriber {
	return &testSubscriber{
		message: message,
	}
}

func (t *testSubscriber) Subscribe() (keybaser.NewSubscription, error) {
	return newTestNewSubscription(t.message), nil
}

type testKeybaseChatAPIClient struct {
	keybaser.KeybaseChatAPIClient
}

func (t *testKeybaseChatAPIClient) ListenForNewTextMessages() (kbchat.NewSubscription, error) {
	return kbchat.NewSubscription{}, nil
}

func (t *testKeybaseChatAPIClient) GetUsername() string {
	return "example"
}

type testResponseWriter struct {
	rw io.ReadWriter
	keybaser.ResponseWriter
}

func (t *testResponseWriter) Reply(text string) {
	_, _ = fmt.Fprint(t.rw, text)
}

func (t *testResponseWriter) ReportError(err error) {
	_, _ = fmt.Fprint(t.rw, err)
}

func testResponseConstructorFunc(rw io.ReadWriter) keybaser.ResponseConstructor {
	return func(_ chat1.ChatChannel, _ keybaser.KeybaseChatAPIClient) keybaser.ResponseWriter {
		return &testResponseWriter{
			rw: rw,
		}
	}
}

func newTestLogger() logrus.FieldLogger {
	l := logrus.New()
	l.SetOutput(new(bytes.Buffer))
	return l
}

func TestKeybaser_Listen(t *testing.T) {
	cl := new(testKeybaseChatAPIClient)

	rw := new(bytes.Buffer)

	message := &kbchat.SubscriptionMessage{
		Message: chat1.MsgSummary{
			Channel: chat1.ChatChannel{
				Name:      "example",
				TopicName: "general",
			},
			Sender: chat1.MsgSender{
				Username: "alice",
			},
			Content: chat1.MsgContent{
				Text: &chat1.MessageText{
					Body: "@example greet alice",
				},
			},
		},
	}

	k, err := keybaser.New(
		cl,
		keybaser.WithLogger(newTestLogger()),
		keybaser.SetResponseConstructor(testResponseConstructorFunc(rw)),
		keybaser.SetSubscriber(newTestSubscriber(message)),
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	k.Command("greet <name>", &keybaser.CommandDefinition{
		Description:       "say greeting",
		Example:           "greet alice",
		AuthorizationFunc: nil,
		Handler: func(request keybaser.Request, response keybaser.ResponseWriter) {
			name := request.Param("name")
			response.Reply("hello " + name)

			cancel()
		},
	})

	var eg errgroup.Group
	eg.Go(func() error {
		err := k.Listen(ctx)
		if err != nil && err != context.Canceled {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	if rw.String() != "hello alice" {
		t.Errorf(`expect "%s" is written but actual "%s"`, "hello alice", rw.String())
	}
}
