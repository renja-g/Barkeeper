package components

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/utils"
)

func SetWinnerComponent() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		winner := "team1"
		if e.ComponentInteraction.Data.CustomID() == "team2_wins_button" {
			winner = "team2"
		}

		matches, err := utils.GetMatches()
		if err != nil {
			return err
		}

		matchID, err := utils.ParseMatchID(e.ComponentInteraction.Message)
		if err != nil {
			return err
		}

		// Update the match with the winner
		for i, match := range matches {
			if match.MatchID == *matchID {
				matches[i].Winner = winner
				break
			}
		}

		err = utils.SaveMatches(matches)
		if err != nil {
			return err
		}

		// Update the participants
		team1Ptr, team2Ptr := utils.ParseTeamMessage(e.Message)

		ratings, err := utils.GetRatings()
		if err != nil {
			return err
		}

		// Update the stats
		for i, rating := range ratings {
			for _, player := range team1Ptr {
				if rating.UserID == *player {
					if winner == "team1" {
						ratings[i].Wins += 1
					} else {
						ratings[i].Losses += 1
					}
					break
				}
			}
			for _, player := range team2Ptr {
				if rating.UserID == *player {
					if winner == "team2" {
						ratings[i].Wins += 1
					} else {
						ratings[i].Losses += 1
					}
					break
				}
			}
		}

		err = utils.SaveRatings(ratings)
		if err != nil {
			return err
		}

		// Update the message with the winner
		winnnerTeam := "Blue"
		if winner == "team2" {
			winnnerTeam = "Red"
		}

		embed := e.ComponentInteraction.Message.Embeds[0]
		embed.Title = "Match Finished"
		embed.Description = winnnerTeam + " wins the match! 🎉"
		embed.Color = 0x00ff00

		return e.UpdateMessage(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			ClearContainerComponents().
			Build(),
		)
	}
}