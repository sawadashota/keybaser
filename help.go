package keybaser

import (
	"strings"
)

// helpCommand shows commands list and usage
func (k *Keybaser) helpCommand() *CommandDefinition {
	return &CommandDefinition{
		Description: "show commands list and usage",
		Handler: func(request Request, response ResponseWriter) {
			var msg strings.Builder
			msg.Grow(len(k.botCommands) * 9)

			for _, cmd := range k.botCommands {
				msg.WriteString(":white_check_mark: *")
				msg.WriteString(cmd.Usage())
				msg.WriteString("* - ")
				msg.WriteString(cmd.Definition().Description)
				if cmd.Definition().AuthorizationFunc != nil {
					msg.WriteString(":lock:")
				}
				msg.WriteString("\n")
				if len(cmd.Definition().Example) > 0 {
					msg.WriteString("> Example: ")
					msg.WriteString(cmd.Definition().Example)
					msg.WriteString("\n")
				}
			}

			response.Reply(msg.String())
		},
	}
}
