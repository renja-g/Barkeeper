package commands

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var help = discord.SlashCommandCreate{
	Name:        "help",
	Description: "Shows the help message.",
}

func HelpHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		var helpMessage strings.Builder
		helpMessage.WriteString("Here are the available commands:\n\n")

		for _, cmd := range Commands {
			switch c := cmd.(type) {
			case discord.SlashCommandCreate:
				helpMessage.WriteString(fmt.Sprintf("**/%s** - %s\n", c.Name, c.Description))
			}
		}

		embed := discord.NewEmbedBuilder().
			SetTitle("Help").
			SetDescription(helpMessage.String()).
			SetColor(0x3498db).
			Build()

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}
}
