Keybaser [![GoDoc](https://godoc.org/github.com/sawadashota/keybaser?status.svg)](https://godoc.org/github.com/sawadashota/keybaser) [![Go Report Card](https://goreportcard.com/badge/github.com/sawadashota/keybaser)](https://goreportcard.com/report/github.com/sawadashota/keybaser) [![GolangCI](https://golangci.com/badges/github.com/sawadashota/keybaser.svg)](https://golangci.com/r/github.com/sawadashota/keybaser)
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
===

Built on top of the Keybase API [github.com/keybase/go-keybase-chat-bot](https://github.com/keybase/go-keybase-chat-bot) with the idea to simplify the Real-Time Messaging feature to easily create Keybase Bots, likely [github.com/shomali11/slacker](https://github.com/shomali11/slacker)

Features
---

Features are almost same as [github.com/shomali11/slacker](https://github.com/shomali11/slacker) 

* Easy definitions of commands and their input
* Available bot initialization, errors and default handlers
* Simple parsing of String, Integer, Float and Boolean parameters
* Contains support for `context.Context`
*Built-in `help` command
* supports authorization
* bot responds to mentions and direct messages
* handlers run concurrently via goroutine
* Full access to the Keybase API [github.com/keybase/go-keybase-chat-bot](https://github.com/keybase/go-keybase-chat-bot)

Dependencies
---

* `commander` [github.com/shomali11/slacker](https://github.com/shomali11/slacker)
* `go-keybase-chat-bot` [github.com/keybase/go-keybase-chat-bot](https://github.com/keybase/go-keybase-chat-bot)

Install
---

```
go get github.com/sawadashota/keybaser
```

Examples
---

### Example 1

Defining a simple command

```go
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
```

Example 2
---

Defining a command with parameter, description and example

```go
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
```