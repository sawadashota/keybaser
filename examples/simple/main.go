package main

import (
	"context"
	"log"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/sawadashota/keybaser"
)

func main() {
	client, err := kbchat.Start(kbchat.RunOptions{
		Oneshot: &kbchat.OneshotOptions{
			Username: "<Bot's Username>",
			PaperKey: "<Bot's PaperKey>",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot, err := keybaser.New(client)
	if err != nil {
		log.Fatal(err)
	}

	definition := &keybaser.CommandDefinition{
		Handler: func(request keybaser.Request, response keybaser.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
