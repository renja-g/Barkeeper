package components

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

func ReshuffleComponent() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		team1, team2 := utils.ParseTeamMessage(e.Message)
		participantIds := append(team1, team2...)
		participants := make([]constants.Rating, 0)

		ratings, err := utils.GetRatings()
		if err != nil {
			return err
		}

		for _, rating := range ratings {
			for _, id := range participantIds {
				if rating.UserID == id {
					ratingCopy := rating
					participants = append(participants, ratingCopy)
					break
				}
			}
		}

		newTeam1, newTeam2 := utils.GenerateTeams(participants)
		newTeam1Rating, newTeam2Rating := utils.CalculateTeamRating(newTeam1), utils.CalculateTeamRating(newTeam2)

		// Create the embed
		embed := discord.NewEmbedBuilder().
			SetTitle("Teams").
			SetColor(0x3498db).
			AddField(fmt.Sprintf("Blue (%d)", newTeam1Rating), utils.FormatTeam(newTeam1), false).
			AddField(fmt.Sprintf("Red (%d)", newTeam2Rating), utils.FormatTeam(newTeam2), false).
			SetFooter("Teams reshuffled", "").
			Build()

		return e.UpdateMessage(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
	}
}
