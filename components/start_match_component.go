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

	// Save the match to the matches.json file
	matches, err := utils.GetMatches()
	if err != nil {
		return err
	}
	matches = append(matches, match)
	err = utils.SaveMatches(matches)
	if err != nil {
		return err
	}

	// Update the message with the match ID
	embed := e.Message.Embeds[0]
	embed.Title = "Match started"
	embed.Footer = &discord.EmbedFooter{
		Text: "MatchID: " + match.MatchID.String(),
	}

	return e.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(embed).
		AddActionRow(
			discord.NewPrimaryButton("Team 1 won", "team1_won_button"),
			discord.NewPrimaryButton("Team 2 won", "team2_won_button"),
			discord.NewDangerButton("Cancel match", "cancel_match_button"),
		).
		Build(),
	)
}
