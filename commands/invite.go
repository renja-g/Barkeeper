package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var invite = discord.SlashCommandCreate{
	Name:        "invite",
	Description: "Sends a message, inviting all online users to play a game.",
}

func InviteHandler(b *dbot.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		ratings, err := utils.GetRatings()
		if err != nil {
			return err
		}

		guildID := e.GuildID()
		onlineUsers := make([]constants.Rating, 0)

		for _, rating := range ratings {
			isOnline := isUserOnline(b, *guildID, rating.UserID)
			if isOnline {
				onlineUsers = append(onlineUsers, rating)
			}
		}

		// Send message where all users get pinged
		message := "Hey, let's play a custom game! "
		for _, rating := range onlineUsers {
			message += fmt.Sprintf("<@%s> ", rating.UserID)
		}

		return e.CreateMessage(discord.MessageCreate{
			Content: message,
		})
	}
}