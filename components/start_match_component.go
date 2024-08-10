package components

import (
	"fmt"
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

		// Move members into their respective voice channels
		if err := moveTeamMembers(e, team1, cfg.BlueChannelID); err != nil {
			log.Printf("Error moving Team Blue members: %v", err)
		}

		if err := moveTeamMembers(e, team2, cfg.RedChannelID); err != nil {
			log.Printf("Error moving Team Red members: %v", err)
		}

		// Update the message with the match ID
		embed := e.Message.Embeds[0]
		embed.Title = "Match started"
		embed.Footer = &discord.EmbedFooter{
			Text: "MatchID: " + match.MatchID.String(),
		}

		// Save the match to the matches.json file
		if err := saveMatch(match); err != nil {
			return fmt.Errorf("failed to save match: %w", err)
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

func moveTeamMembers(e *handler.ComponentEvent, team []snowflake.ID, channelID snowflake.ID) error {
	for _, memberID := range team {
		_, err := e.Client().Rest().UpdateMember(*e.GuildID(), memberID, discord.MemberUpdate{
			ChannelID: &channelID,
		})
		if err != nil {
			log.Printf("Error moving member %s to channel %s: %v", memberID, channelID, err)
		}
	}
	return nil
}

func saveMatch(match constants.Match) error {
	matches, err := utils.GetMatches()
	if err != nil {
		return fmt.Errorf("failed to get matches: %w", err)
	}
	matches = append(matches, match)
	if err := utils.SaveMatches(matches); err != nil {
		return fmt.Errorf("failed to save matches: %w", err)
	}
	return nil
}
