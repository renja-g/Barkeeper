package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var info = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Displays information about a given user.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to get information about.",
			Required:    true,
		},
	},
}

func InfoHandler(e *handler.CommandEvent) error {
	ratings, err := utils.GetRatings()
	if err != nil {
		return err
	}

	// Check if the user has a rating
	userRating := constants.Rating{}
	for _, rating := range ratings {
		if rating.UserID == e.SlashCommandInteractionData().User("user").ID {
			userRating = rating
			break
		}
	}

	if userRating.UserID == 0 {
		embed := discord.NewEmbedBuilder().
			SetTitle("User not found").
			SetDescriptionf("User %s not found.", e.SlashCommandInteractionData().User("user").Mention()).
			SetColor(0xff0000).
			Build()

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}

	winrate := 0.0
	if userRating.Wins+userRating.Losses > 0 {
		winrate = float64(userRating.Wins) / float64(userRating.Wins+userRating.Losses) * 100
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("User info").
		SetDescriptionf("User: %s\nRating: %d\nWins: %d\nLosses: %d\nWinrate: %.2f%%", e.SlashCommandInteractionData().User("user").Mention(), userRating.Rating, userRating.Wins, userRating.Losses, winrate).
		SetColor(0x3498db).
		Build()

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}
