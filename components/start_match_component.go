package components

import (
	"log"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

func StartMatchComponent(cfg *dbot.Config) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		team1Ptr, team2Ptr := utils.ParseTeamMessage(e.Message)

		// Dereference the pointers in the slices
		team1 := make([]snowflake.ID, len(team1Ptr))
		copy(team1, team1Ptr)

		team2 := make([]snowflake.ID, len(team2Ptr))
		copy(team2, team2Ptr)

		match := constants.Match{
			MatchID:   snowflake.New(time.Now()),
			Team1:     team1,
			Team2:     team2,
			Timestamp: time.Now().Unix(),
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

		// Move member into their respective voice channels
		blueChannelID := cfg.BlueChannelID
		redChannelID := cfg.RedChannelID

		for _, memberID := range team1 {
			_, err := e.Client().Rest().UpdateMember(*e.GuildID(), memberID, discord.MemberUpdate{
				ChannelID: &blueChannelID,
			})
			if err != nil {
				log.Printf("error moving member %s to blue channel: %s", memberID, err)
				continue
			}
		}

		for _, memberID := range team2 {
			_, err := e.Client().Rest().UpdateMember(*e.GuildID(), memberID, discord.MemberUpdate{
				ChannelID: &redChannelID,
			})
			if err != nil {
				log.Printf("error moving member %s to red channel: %s", memberID, err)
				continue
			}
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
				discord.NewPrimaryButton("Team Blue wins", "team1_wins_button"),
				discord.NewPrimaryButton("Team Red wins", "team2_wins_button"),
				discord.NewDangerButton("Cancel match", "cancel_match_button"),
			).
			Build(),
		)
	}
}
