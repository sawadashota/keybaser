package keybaser_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/sawadashota/keybaser"
)

func TestBotCommand(t *testing.T) {
	// To test CommandDefinition.Handler
	buffer := new(bytes.Buffer)

	cmd := keybaser.NewBotCommand("greet <name>", &keybaser.CommandDefinition{
		Description:       "Just reply greeting",
		Example:           "greet John",
		AuthorizationFunc: nil,
		Handler: func(request keybaser.Request, response keybaser.ResponseWriter) {
			name := request.StringParam("name", "")
			buffer.Write([]byte("Hello " + name))
		},
	})

	t.Run("Usage", func(t *testing.T) {
		if cmd.Usage() != "greet <name>" {
			t.Errorf("Usage() is expected to return %s but actual %s", "greet <name>", cmd.Usage())
		}
	})

	t.Run("Definition", func(t *testing.T) {
		if cmd.Definition().Example != "greet John" {
			t.Errorf("Usage() is expected to return %s but actual %s", "greet John", cmd.Definition().Example)
		}
	})

	t.Run("Match", func(t *testing.T) {
		parameters, isMatch := cmd.Match("@example greet John")
		if isMatch != true {
			t.Errorf("Match() expected to be true but actual false")
		}

		if parameters.StringParam("name", "") != "John" {
			t.Errorf("parameters is expected to have %s but actual %s", "John", parameters.StringParam("name", ""))
		}
	})

	t.Run("Match", func(t *testing.T) {
		parameters, isMatch := cmd.Match("@example greet John")
		if isMatch != true {
			t.Errorf("Match() expected to be true but actual false")
		}

		if parameters.StringParam("name", "") != "John" {
			t.Errorf("parameters is expected to have %s but actual %s", "John", parameters.StringParam("name", ""))
		}
	})

	t.Run("Execute", func(t *testing.T) {
		parameters, isMatch := cmd.Match("@example greet John")

		if !isMatch {
			t.Fatal("command should be match")
		}

		if buffer.String() != "" {
			t.Fatal("buffer should be empty before executing")
		}

		var message kbchat.SubscriptionMessage
		req := keybaser.NewRequest(context.Background(), &message, parameters)
		resp := keybaser.NewResponse(message.Message.Channel, nil)
		cmd.Execute(req, resp)

		if buffer.String() != "Hello John" {
			t.Errorf("buffer is expect to be %s but actual %s", "Hello John", buffer.String())
		}
	})
}
