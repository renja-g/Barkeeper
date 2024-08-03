package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	dbot "github.com/renja-g/Barkeeper"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
)

// RiotAccountResponse represents the structure of the Riot API response
type RiotAccountResponse struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// AccountData extends RiotAccountResponse with additional information
type AccountData struct {
	RiotAccountResponse
	Region        string
	VerifyImageID int
}

var DataCache = make(map[string]AccountData)

var link_account = discord.SlashCommandCreate{
	Name:        "link_account",
	Description: "Link a League of Legends account to your Discord account.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "game_name",
			Description: "Your League of Legends in-game name.",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "tag_line",
			Description: "Your League of Legends tag line.",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "region",
			Description: "Your League of Legends region.",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "BR (Brazil)", Value: "BR1"},
				{Name: "EUNE (Europe Nordic & East)", Value: "EUN1"},
				{Name: "EUW (Europe West)", Value: "EUW1"},
				{Name: "JP (Japan)", Value: "JP1"},
				{Name: "KR (Korea)", Value: "KR"},
				{Name: "LAN (Latin America North)", Value: "LA1"},
				{Name: "LAS (Latin America South)", Value: "LA2"},
				{Name: "NA (North America)", Value: "NA1"},
				{Name: "OCE (Oceania)", Value: "OC1"},
				{Name: "TR (Turkey)", Value: "TR1"},
				{Name: "RU (Russia)", Value: "RU"},
				{Name: "PH (Philippines)", Value: "PH2"},
				{Name: "SG (Singapore)", Value: "SG2"},
				{Name: "TH (Thailand)", Value: "TH2"},
				{Name: "TW (Taiwan)", Value: "TW2"},
				{Name: "VN (Vietnam)", Value: "VN2"},
			},
		},
	},
}

func LinkAccountHandler(cfg *dbot.Config) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		gameName := data.String("game_name")
		tagLine := data.String("tag_line")
		region := data.String("region")

		url := fmt.Sprintf("https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", gameName, tagLine, cfg.RiotApiKey)
		resp, err := http.Get(url)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to verify account. Please try again later.").Build())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to read response. Please try again later.").Build())
			}

			var riotResponse RiotAccountResponse
			err = json.Unmarshal(body, &riotResponse)
			if err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to parse response. Please try again later.").Build())
			}

			// Generate a random image ID the user will have to change his profile picture to
			randomIconID := rand.Intn(29)
			imageURL := fmt.Sprintf("https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default//v1/profile-icons/%d.jpg", randomIconID)

			accountData := AccountData{
				RiotAccountResponse: riotResponse,
				Region:              region,
				VerifyImageID:       randomIconID,
			}
			dataID := uuid.New().String()
			DataCache[dataID] = accountData

			embed := discord.NewEmbedBuilder().
				SetTitle("Verify Account").
				SetDescription(fmt.Sprintf("Account %s#%s has been found.\nChange you profile picture to the image below and click the verify button.", accountData.GameName, accountData.TagLine)).
				SetImage(imageURL).
				Build()

			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetEmbeds(embed).
				AddActionRow(
					discord.NewPrimaryButton("Verify", "verify_acc/"+dataID),
				).
				Build(),
			)
		} else {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Account not found or API error. Please check your details and try again.").Build())
		}
	}
}
