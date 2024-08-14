package components

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/utils"
)

func SetWinnerComponent(cfg *dbot.Config) handler.ButtonComponentHandler {
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
			if match.MatchID == matchID {
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

		profiles, err := utils.GetProfiles()
		if err != nil {
			return err
		}

		// Update the stats
		for i, profile := range profiles {
			for _, player := range team1Ptr {
				if profile.UserID == player {
					if winner == "team1" {
						profiles[i].Wins += 1
					} else {
						profiles[i].Losses += 1
					}
					break
				}
			}
			for _, player := range team2Ptr {
				if profile.UserID == player {
					if winner == "team2" {
						profiles[i].Wins += 1
					} else {
						profiles[i].Losses += 1
					}
					break
				}
			}
		}

		err = utils.SaveProfiles(profiles)
		if err != nil {
			return err
		}

		// Move members back to the lobby
		participants := append(team1Ptr, team2Ptr...)
		moveTeamMembers(e, participants, cfg.LobbyChannelID)

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
