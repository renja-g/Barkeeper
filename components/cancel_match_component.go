package components

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

func CancelMatchComponent(cfg *dbot.Config) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		// Check if the user has the admin role
		member, err := e.Client().Rest().GetMember(*e.GuildID(), e.User().ID)
		if err != nil {
			return err
		}

		if !utils.HasAdminRole(member, cfg.AdminRoleID) {
			return nil
		}

		if err := e.DeferUpdateMessage(); err != nil {
			return fmt.Errorf("failed to defer update: %w", err)
		}

		matches, err := utils.GetMatches()
		if err != nil {
			return err
		}

		matchID, err := utils.ParseMatchID(e.Message)
		if err != nil {
			return err
		}

		// Remove the match
		newMatchArr := make([]constants.Match, 0)
		for i, match := range matches {
			if match.MatchID != matchID {
				newMatchArr = append(newMatchArr, matches[i])
			}
		}

		err = utils.SaveMatches(newMatchArr)
		if err != nil {
			return err
		}

		// Move members back to the lobby
		team1Ptr, team2Ptr := utils.ParseTeamMessage(e.Message)
		participants := append(team1Ptr, team2Ptr...)

		moveTeamMembers(e, participants, cfg.LobbyChannelID)

		embed := e.Message.Embeds[0]
		embed.Title = "Match Cancelled"
		embed.Color = 0xff0000

		_, err = e.Client().Rest().UpdateInteractionResponse(
			e.ApplicationID(),
			e.Token(),
			discord.NewMessageUpdateBuilder().
				SetEmbeds(embed).
				ClearContainerComponents().
				Build(),
		)
		if err != nil {
			return fmt.Errorf("failed to update interaction response: %w", err)
		}

		return nil
	}
}
