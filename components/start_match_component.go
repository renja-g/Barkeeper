package components

import (
	"github.com/bwmarrin/snowflake"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	dSnowflake "github.com/disgoorg/snowflake/v2"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

func StartMatchComponent(e *handler.ComponentEvent) error {
    team1Ptr, team2Ptr := utils.ParseTeamMessage(e.Message)

    // Dereference the pointers in the slices
    team1 := make([]dSnowflake.ID, len(team1Ptr))
    for i, v := range team1Ptr {
        team1[i] = *v
    }

    team2 := make([]dSnowflake.ID, len(team2Ptr))
    for i, v := range team2Ptr {
        team2[i] = *v
    }

    Node, err := snowflake.NewNode(1)
    if err != nil {
        return err
    }

    match := constants.Match{
        MatchID: Node.Generate(),
        Team1:   team1,
        Team2:   team2,
    }


	return e.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	})
}
