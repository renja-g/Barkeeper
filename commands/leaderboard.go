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
		ratings, err := utils.GetRatings()
		if err != nil {
			return err
		}

		// Filter out users with less than 5 games
		ratings = filterRatings(ratings, 5)

		// Sort ratings and get top 15

		ratings = utils.SortRatingsByWinRate(ratings)
		if len(ratings) > 15 {
			ratings = ratings[:15]
		}

		fields := make([]discord.EmbedField, len(ratings))
		inline := true
		for i, r := range ratings {
			winrate := 0.0
			if r.Wins+r.Losses > 0 {
				winrate = float64(r.Wins) / float64(r.Wins+r.Losses) * 100
			}

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

func filterRatings(ratings []constants.Rating, minGames int) []constants.Rating {
	filtered := make([]constants.Rating, 0)
	for _, r := range ratings {
		if r.Wins+r.Losses >= minGames {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
