package keybaser

import (
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Description       string
	Example           string
	AuthorizationFunc func(request Request) bool
	Handler           func(request Request, response ResponseWriter)
}

// NewBotCommand creates a new bot command object
func NewBotCommand(usage string, definition *CommandDefinition) BotCommand {
	command := commander.NewCommand(usage)
	return &botCommand{
		usage:      usage,
		definition: definition,
		command:    command,
	}
}

// botCommand structure contains the bot's command, description and handler
type botCommand struct {
	usage      string
	definition *CommandDefinition
	command    *commander.Command
}

// BotCommand interface
type BotCommand interface {
	Usage() string
	Definition() *CommandDefinition

	Match(text string) (*proper.Properties, bool)
	Execute(request Request, response ResponseWriter)
}

// Usage returns the command usage
func (c *botCommand) Usage() string {
	return c.usage
}

// Description returns the command description
func (c *botCommand) Definition() *CommandDefinition {
	return c.definition
}

// Match determines whether the bot should respond based on the text received
func (c *botCommand) Match(text string) (*proper.Properties, bool) {
	return c.command.Match(text)
}

// Execute executes the handler logic
func (c *botCommand) Execute(request Request, response ResponseWriter) {
	if c.definition == nil || c.definition.Handler == nil {
		return
	}
	c.definition.Handler(request, response)
}
