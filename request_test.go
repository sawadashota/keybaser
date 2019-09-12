package keybaser_test

import (
	"context"
	"testing"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/sawadashota/keybaser"
	"github.com/shomali11/proper"
)

func newTestingSubscriptionMessage(t *testing.T) *kbchat.SubscriptionMessage {
	t.Helper()

	return &kbchat.SubscriptionMessage{
		Message: chat1.MsgSummary{
			Content: chat1.MsgContent{
				TypeName: "text",
				Text: &chat1.MessageText{
					Body: "@example hello",
				},
			},
		},
	}
}

func newTestingProperties(t *testing.T) *proper.Properties {
	t.Helper()

	params := map[string]string{
		"word": "hello",
		"ok":   "true",
		"age":  "10",
		"pi":   "3.14",
	}
	return proper.NewProperties(params)
}

func TestRequest(t *testing.T) {
	ctx := context.Background()
	message := newTestingSubscriptionMessage(t)
	properties := newTestingProperties(t)
	req := keybaser.NewRequest(ctx, message, properties)

	t.Run("Param", func(t *testing.T) {
		if req.Param("word") != "hello" {
			t.Errorf("req.Param() expected to return %s, but actual %s", "hello", req.Param("word"))
		}
	})

	t.Run("StringParam", func(t *testing.T) {
		if req.StringParam("word", "hi") != "hello" {
			t.Errorf("req.StringParam() expected to return %s, but actual %s", "hello", req.StringParam("word", "hi"))
		}
		if req.StringParam("hoge", "hi") != "hi" {
			t.Errorf("req.StringParam() expected to return %s, but actual %s", "hi", req.StringParam("hoge", "hi"))
		}
	})

	t.Run("BooleanParam", func(t *testing.T) {
		if req.BooleanParam("ok", false) != true {
			t.Errorf("req.BooleanParam() is expected to return %t, but actual %t", true, req.BooleanParam("ok", false))
		}
		if req.BooleanParam("hoge", true) != true {
			t.Errorf("req.BooleanParam() is expected to return %t, but actual %t", true, req.BooleanParam("hoge", true))
		}
	})

	t.Run("IntegerParam", func(t *testing.T) {
		if req.IntegerParam("age", 100) != 10 {
			t.Errorf("req.IntegerParam() is expected to return %d, but actual %d", 10, req.IntegerParam("age", 100))
		}
		if req.IntegerParam("hoge", 100) != 100 {
			t.Errorf("req.IntegerParam() is expected to return %d, but actual %d", 100, req.IntegerParam("hoge", 100))
		}
	})

	t.Run("FloatParam", func(t *testing.T) {
		if req.FloatParam("pi", 100.111) != 3.14 {
			t.Errorf("req.FloatParam() is expected to return %f, but actual %f", 3.14, req.FloatParam("pi", 100.111))
		}
		if req.FloatParam("hoge", 100.111) != 100.111 {
			t.Errorf("req.FloatParam() is expected to return %f, but actual %f", 100.111, req.FloatParam("hoge", 100.111))
		}
	})

	t.Run("Message", func(t *testing.T) {
		if req.Message() != message {
			t.Errorf("req.Message() is expected to return %v, but actual %v", message, req.Message())
		}
	})

	t.Run("Properties", func(t *testing.T) {
		if req.Properties() != properties {
			t.Errorf("req.Properties() is expected to return %v, but actual %v", properties, req.Properties())
		}
	})
}
