package components

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

func CancelMatchComponent() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
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

		embed := e.Message.Embeds[0]
		embed.Title = "Match Cancelled"
		embed.Color = 0xff0000

		return e.UpdateMessage(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			ClearContainerComponents().
			Build(),
		)
	}
}
