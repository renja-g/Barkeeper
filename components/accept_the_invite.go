package components

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func AcceptTheInviteComponent() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		// Add or remove the user that pressed the button to/from the embed
		oldEmbed := e.ComponentInteraction.Message.Embeds[0]
		userID := fmt.Sprintf("<@%s>", e.ComponentInteraction.User().ID)

		// Check if the user is already in the fields
		var newFields []discord.EmbedField
		userExists := false
		for _, field := range oldEmbed.Fields {
			if field.Value == userID {
				userExists = true
			} else {
				newFields = append(newFields, field)
			}
		}

		// If user does not exist, add them
		if !userExists {
			inline := true
			newFields = append(newFields, discord.EmbedField{
				Name:   "",
				Value:  userID,
				Inline: &inline,
			})
		}

		// Update the description with the new count
		count := len(newFields)
		description := fmt.Sprintf("(%d/10)", count)
		if oldEmbed.Description != "" {
			oldCount := count
			fmt.Sscanf(oldEmbed.Description, "(%d/10)", &oldCount)
			description = strings.Replace(oldEmbed.Description, fmt.Sprintf("(%d/10)", oldCount), description, 1)
		}

		newEmbed := discord.NewEmbedBuilder().
			SetTitle(oldEmbed.Title).
			SetDescription(description).
			SetFields(newFields...).
			Build()

		return e.UpdateMessage(discord.NewMessageUpdateBuilder().
			SetEmbeds(newEmbed).
			Build(),
		)
	}
}
