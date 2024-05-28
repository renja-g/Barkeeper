package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/utils"
)

var history = discord.SlashCommandCreate{
	Name:        "history",
	Description: "Show the history of all matches.",
}

func HistoryHandler(e *handler.CommandEvent) error {
	matches, err := utils.GetMatches()
	if err != nil {
		return err
	}


	fields := make([]discord.EmbedField, len(matches))
	inline := true
	for i, m := range matches {
		if m.Winner == "" {
			continue
		}
		
		winner := "Blue"
		if m.Winner == "team2" {
			winner = "Red"
		}

		team1String := ""
		for _, id := range m.Team1 {
			team1String += fmt.Sprintf("<@%s>\n", id)
		}

		team2String := ""
		for _, id := range m.Team2 {
			team2String += fmt.Sprintf("<@%s>\n", id)
		}

		fields[i] = discord.EmbedField{
			Name:   fmt.Sprintf("Team **%s** won:", winner),
			Value:  fmt.Sprintf("Blue\n%s\nRed\n%s\n<t:%d:R>", team1String, team2String, m.Timestamp),
			Inline: &inline,
		}
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Leaderboard").
		SetColor(0x3498db).
		SetFields(fields...).
		Build()

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}
