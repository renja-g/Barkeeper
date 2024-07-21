package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/utils"
)

var leaderboard = discord.SlashCommandCreate{
	Name:        "leaderboard",
	Description: "Displays the leaderboard.",
}

func LeaderboardHandler(e *handler.CommandEvent) error {
	ratings, err := utils.GetRatings()
	if err != nil {
		return err
	}

	ratings = utils.GetLeaderboard(ratings)
	if len(ratings) > 20 {
		ratings = ratings[:20]
	}

	fields := make([]discord.EmbedField, len(ratings))
	inline := true
	for i, r := range ratings {
		winrate := 0.0
		if r.Wins+r.Losses > 0 {
			winrate = float64(r.Wins) / float64(r.Wins+r.Losses) * 100
		}

		fields[i] = discord.EmbedField{
			Value:  fmt.Sprintf("<@%s>\nRating: %d\nWins: %d\nLosses: %d\nWinrate: %.2f%%", r.UserID, r.Rating, r.Wins, r.Losses, winrate),
			Name:   "Rank",
			Inline: &inline,
		}
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Leaderboard").
		SetColor(0x3498db).
		SetFields(fields...).
		SetFooterText("Top 20 players").
		Build()

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}
