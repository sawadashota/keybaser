package main

import (
	"context"
	"log"
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/sawadashota/keybaser"
)

func main() {
	username := os.Getenv("KEYBASE_USERNAME")
	paperkey := os.Getenv("KEYBASE_PAPERKEY")

	client, err := kbchat.Start(kbchat.RunOptions{
		Oneshot: &kbchat.OneshotOptions{
			Username: username,
			PaperKey: paperkey,
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
		Description: "Greet to you",
		Example:     "greet alice",
		Handler: func(request keybaser.Request, response keybaser.ResponseWriter) {
			name := request.Param("name")
			response.Reply("Hello " + name)
		},
	}

	bot.Command("greet <name>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
