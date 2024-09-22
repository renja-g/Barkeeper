package commands

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var teams = discord.SlashCommandCreate{
	Name:        "teams",
	Description: "Generates fair teams with the users in the voice channel you are in.",
}

func TeamsHandler(b *dbot.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		// Check if the user is in a voice channel
		voiceState, ok := b.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
		if !ok {
			embed := discord.NewEmbedBuilder().
				SetTitle("User not found").
				SetDescriptionf("You are not in a voice channel.").
				SetColor(0xff0000).
				Build()

			return e.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		}

		// Get the voice channel
		audioChannel, ok := b.Client.Caches().GuildAudioChannel(*voiceState.ChannelID)
		if !ok {
			embed := discord.NewEmbedBuilder().
				SetTitle("Channel not found").
				SetDescriptionf("The voice channel you are in was not found.").
				SetColor(0xff0000).
				Build()

			return e.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		}

		// Get the members in the voice channel
		var ids []snowflake.ID
		b.Client.Caches().VoiceStatesForEach(audioChannel.GuildID(), func(state discord.VoiceState) {
			if state.ChannelID != nil && *state.ChannelID == audioChannel.ID() {
				ids = append(ids, state.UserID)
			}
		})

		// Check if there are exactly 10 users in the voice channel
		if len(ids) != 10 {
			embed := discord.NewEmbedBuilder().
				SetTitle("Invalid team").
				SetDescriptionf("The voice channel you are in does not have exactly 10 members.").
				SetColor(0xff0000).
				Build()

			return e.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		}

		// Check if all members have a profile
		profiles, err := utils.GetProfiles()
		if err != nil {
			return err
		}

		var missingProfiles []snowflake.ID
		for _, id := range ids {
			found := false
			for _, profile := range profiles {
				if profile.UserID == id {
					found = true
					break
				}
			}
			if !found {
				missingProfiles = append(missingProfiles, id)
			}
		}

		// If there are missing profiles, return an error mentioning the users that are missing a profile
		if len(missingProfiles) > 0 {
			var mentions []string
			for _, id := range missingProfiles {
				mentions = append(mentions, "<@"+id.String()+">")
			}

			mentionString := strings.Join(mentions, ", ")

			embed := discord.NewEmbedBuilder().
				SetTitle("Missing Profiles").
				SetDescriptionf("The following users are missing a profile: %s", mentionString).
				SetColor(0xff0000).
				Build()

			return e.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		}

		// Get the profiles for the members in the voice channel
		var memberProfiles []constants.Profile
		for _, id := range ids {
			for _, profile := range profiles {
				if profile.UserID == id {
					memberProfiles = append(memberProfiles, profile)
					break
				}
			}
		}

		// Generate the best teams
		team1, team2 := utils.GenerateTeams(memberProfiles)
		team1Rating, team2Rating := utils.CalculateTeamRating(team1), utils.CalculateTeamRating(team2)

		// Create the embed
		embed := discord.NewEmbedBuilder().
			SetTitle("Teams").
			SetColor(0x3498db).
			AddField(fmt.Sprintf("Blue (%d)", team1Rating), utils.FormatTeam(team1), false).
			AddField(fmt.Sprintf("Red (%d)", team2Rating), utils.FormatTeam(team2), false).
			Build()

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(embed).
			AddActionRow(
				discord.NewPrimaryButton("Start match", "start_match_button"),
				discord.NewPrimaryButton("Reshuffle", "reshuffle_button")).
			Build(),
		)
	}
}
