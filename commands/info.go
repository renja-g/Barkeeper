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

func InfoHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		profiles, err := utils.GetProfiles()
		if err != nil {
			return err
		}

		// Check if the user has a profile
		userProfile := constants.Profile{}
		for _, profile := range profiles {
			if profile.UserID == e.SlashCommandInteractionData().User("user").ID {
				userProfile = profile
				break
			}
		}

		if userProfile.UserID == 0 {
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
		total := userProfile.Wins + userProfile.Losses
		if total > 0 {
			winrate = float64(userProfile.Wins) / float64(total) * 100
		}

		embed := discord.NewEmbedBuilder().
			SetTitle("User info").
			SetDescriptionf("User: %s\nRating: %d\nWins: %d\nLosses: %d\nWinrate: %.2f%%", e.SlashCommandInteractionData().User("user").Mention(), userProfile.Rating, userProfile.Wins, userProfile.Losses, winrate).
			SetColor(0x3498db).
			Build()

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}
}
