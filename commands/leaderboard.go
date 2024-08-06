package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var leaderboard = discord.SlashCommandCreate{
	Name:        "leaderboard",
	Description: "Displays the leaderboard.",
}

func LeaderboardHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		profiles, err := utils.GetProfiles()
		if err != nil {
			return err
		}

		// Filter out users with less than 5 games
		profiles = filterProfiles(profiles, 5)

		// Sort profiles and get top 15

		profiles = utils.SortProfilesByWinRate(profiles)
		if len(profiles) > 15 {
			profiles = profiles[:15]
		}

		fields := make([]discord.EmbedField, len(profiles))
		for i, r := range profiles {
			winrate := 0.0
			if total := r.Wins + r.Losses; total > 0 {
				winrate = float64(r.Wins) / float64(total) * 100
			}

			inline := true
			winLoss := fmt.Sprintf("%d/%d", r.Wins, r.Losses)
			fields[i] = discord.EmbedField{
				Value:  fmt.Sprintf("<@%s>\nRating: %d\nW/L: %s\nWinrate: %.2f%%", r.UserID, r.Rating, winLoss, winrate),
				Name:   fmt.Sprintf("#%d", i+1),
				Inline: &inline,
			}
		}

		embed := discord.NewEmbedBuilder().
			SetTitle("Leaderboard").
			SetColor(0x3498db).
			SetFields(fields...).
			SetDescription("Top 15 players by winrate with at least 5 games played.").
			Build()

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}
}

func filterProfiles(profiles []constants.Profile, minGames int) []constants.Profile {
	filtered := make([]constants.Profile, 0)
	for _, r := range profiles {
		if r.Wins+r.Losses >= minGames {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
