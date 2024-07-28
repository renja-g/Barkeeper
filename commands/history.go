package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var history = discord.SlashCommandCreate{
	Name:        "history",
	Description: "Show the history of all matches.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "To show the history of a specific user.",
			Required:    false,
		},
	},
}

func HistoryHandler(b *dbot.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		userOption, passed := data.OptUser("user")

		matches, err := utils.GetMatches()
		if err != nil {
			return err
		}

		// Filter matches if a user is specified
		if passed {
			matches = filterMatchesByUser(matches, userOption.ID)
		}

		// reverse the matches so the most recent match is shown first
		for i, j := 0, len(matches)-1; i < j; i, j = i+1, j-1 {
			matches[i], matches[j] = matches[j], matches[i]
		}

		const maxFieldsPerPage = 6
		var pages []*discord.EmbedBuilder

		for i := 0; i < len(matches); i += maxFieldsPerPage {
			end := i + maxFieldsPerPage
			if end > len(matches) {
				end = len(matches)
			}

			pageMatches := matches[i:end]
			embed := createHistoryEmbed(pageMatches, &userOption)
			pages = append(pages, embed)
		}

		return b.Paginator.Create(e.Respond, paginator.Pages{
			ID: e.ID().String(),
			PageFunc: func(page int, embed *discord.EmbedBuilder) {
				*embed = *pages[page]
				embed.SetFooter(fmt.Sprintf("Page %d of %d", page+1, len(pages)), "")
			},
			Pages:      len(pages),
			ExpireMode: paginator.ExpireModeAfterLastUsage,
		}, false)
	}
}
func filterMatchesByUser(matches []constants.Match, userID snowflake.ID) []constants.Match {
	filtered := []constants.Match{}
	for _, match := range matches {
		if containsUser(match.Team1, userID) || containsUser(match.Team2, userID) {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

func containsUser(team []snowflake.ID, userID snowflake.ID) bool {
	for _, id := range team {
		if id == userID {
			return true
		}
	}
	return false
}

func createHistoryEmbed(matches []constants.Match, user *discord.User) *discord.EmbedBuilder {
	title := "Match History"
	if user != nil {
		title = fmt.Sprintf("Match History for %s", user.Username)
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(title).
		SetColor(0x3498db)

	for _, m := range matches {
		if m.Winner == "" {
			continue
		}

		winner := ":blue_square:"
		if m.Winner == "team2" {
			winner = ":red_square:"
		}

		team1String := ""
		for _, id := range m.Team1 {
			team1String += fmt.Sprintf("<@%s>\n", id)
		}

		team2String := ""
		for _, id := range m.Team2 {
			team2String += fmt.Sprintf("<@%s>\n", id)
		}

		fieldName := fmt.Sprintf("Winner: **%s**", winner)
		fieldValue := fmt.Sprintf("Blue\n%s\nRed\n%s\n<t:%d:R>", team1String, team2String, m.Timestamp)

		embed.AddField(fieldName, fieldValue, true)
	}

	return embed
}
