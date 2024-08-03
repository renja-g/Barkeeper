package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/renja-g/Barkeeper/commands"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type SummonerResponse struct {
	ProfileIconId int `json:"profileIconId"`
}

func VerifyAccountLinkComponent() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		dataID := e.Vars["data"]

		accountData, ok := commands.DataCache[dataID]
		if !ok {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Verification data not found. Please try again.").Build())
		}

		// Fetch current summoner data
		url := fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s?api_key=RGAPI-55cd5b47-9656-4ebe-ac35-086b704432f4", accountData.PUUID)
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

		var summonerResponse SummonerResponse
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
				SetDescription(fmt.Sprintf("Account %s#%s has been successfully verified for region %s.", accountData.GameName, accountData.TagLine, accountData.Region)).
				SetColor(0x00FF00) // Green color

			// TODO, save the verified account data

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
