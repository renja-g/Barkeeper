package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/commands"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func VerifyAccountLinkComponent(cfg *dbot.Config) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		dataID := e.Vars["data"]

		accountData, ok := commands.DataCache[dataID]
		if !ok {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Verification data not found. Please try again.").Build())
		}

		// Fetch current summoner data
		url := fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s?api_key=%s", accountData.PUUID, cfg.RiotApiKey)
		resp, err := http.Get(url)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to verify account. Please try again later.").Build())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to fetch summoner data. Please try again later.").Build())
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to read response. Please try again later.").Build())
		}

		var summonerResponse constants.SummonerResponse
		err = json.Unmarshal(body, &summonerResponse)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to parse response. Please try again later.").Build())
		}

		oldEmbed := e.ComponentInteraction.Message.Embeds[0]
		var newEmbed *discord.EmbedBuilder

		success := false
		if summonerResponse.ProfileIconId == accountData.VerifyImageID {
			// Verification successful
			success = true
			newEmbed = discord.NewEmbedBuilder().
				SetTitle("Account Verified").
				SetDescription(fmt.Sprintf("Account %s#%s has been successfully linked to your account.", accountData.GameName, accountData.TagLine)).
				SetColor(0x00FF00) // Green color

			userId := e.User().ID
			profiles, err := utils.GetProfiles()
			if err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to fetch profiles. Please try again later.").Build())
			}

			// Check if the user already has a profile
			var userProfile constants.Profile
			for _, profile := range profiles {
				if profile.UserID == userId {
					userProfile = profile
					break
				}
			}
			if userProfile.UserID == 0 {
				// Create a new profile
				userProfile = constants.Profile{
					UserID:        userId,
					Rating:        0,
					Wins:          0,
					Losses:        0,
					VerifiedPUUID: &accountData.PUUID,
				}
				profiles = append(profiles, userProfile)
			} else {
				// Update the existing profile
				userProfile.VerifiedPUUID = &accountData.PUUID
				for i, profile := range profiles {
					if profile.UserID == userId {
						profiles[i] = userProfile
						break
					}
				}
			}

			err = utils.SaveProfiles(profiles)
			if err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to save profile. Please try again later.").Build())
			}
		} else {
			// Verification failed
			success = false
			newEmbed = discord.NewEmbedBuilder().
				SetTitle("Verification Failed").
				SetDescription("The profile icon doesn't match. Please set your icon to the one shown and try again.").
				SetImage(oldEmbed.Image.URL). // Keep showing the required icon
				SetColor(0xFF0000)            // Red color
		}

		if success {
			delete(commands.DataCache, dataID)
			return e.UpdateMessage(discord.NewMessageUpdateBuilder().
				SetEmbeds(newEmbed.Build()).
				ClearContainerComponents().
				Build(),
			)
		} else {
			return e.UpdateMessage(discord.NewMessageUpdateBuilder().
				SetEmbeds(newEmbed.Build()).
				AddActionRow(
					discord.NewPrimaryButton("Verify", "verify_acc/"+dataID),
				).
				Build(),
			)
		}
	}
}
