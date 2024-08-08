package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
)

var DataCache = make(map[string]constants.AccountData)

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
			Choices:     getRegionChoices(),
		},
	},
}

func getRegionChoices() []discord.ApplicationCommandOptionChoiceString {
	return []discord.ApplicationCommandOptionChoiceString{
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
	}
}

func LinkAccountHandler(cfg *dbot.Config) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if err := e.DeferCreateMessage(true); err != nil {
			return fmt.Errorf("failed to defer message: %w", err)
		}

		gameName := data.String("game_name")
		tagLine := data.String("tag_line")
		region := data.String("region")

		riotResponse, err := verifyRiotAccount(cfg.RiotApiKey, gameName, tagLine)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(err.Error()).Build())
		}

		summonerResponse, err := fetchSummonerData(cfg.RiotApiKey, riotResponse.PUUID, region)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(err.Error()).Build())
		}

		accountData := constants.AccountData{
			RiotAccountResponse: *riotResponse,
			Region:              region,
			VerifyImageID:       getRandomIcon(summonerResponse.ProfileIconId),
		}

		dataID := uuid.New().String()
		DataCache[dataID] = accountData

		return sendVerificationMessage(e, accountData, dataID)
	}
}

func verifyRiotAccount(apiKey, gameName, tagLine string) (*constants.RiotAccountResponse, error) {
	url := fmt.Sprintf("https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", gameName, tagLine, apiKey)
	return makeRequest[constants.RiotAccountResponse](url, "Failed to verify account")
}

func fetchSummonerData(apiKey, puuid, region string) (*constants.SummonerResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s?api_key=%s", region, puuid, apiKey)
	return makeRequest[constants.SummonerResponse](url, "Failed to fetch summoner data")
}

func makeRequest[T any](url, errorMessage string) (*T, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorMessage, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: status code %d", errorMessage, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

var icons = map[int]string{
	7:  "DEBONAIR ROSE ICON",
	9:  "DAGGERS ICON",
	10: "WINGED SWORD ICON",
	12: "FULLY STACKED MEJAI'S ICON",
	18: "MIX MIX ICON",
	21: "TREE OF LIFE ICON",
	22: "REVIVE ICON",
	23: "LIL' SPROUT ICON",
	24: "SPIKE SHIELD ICON",
	28: "TIBBERS ICON",
}

func getRandomIcon(currentIconId int) int {
	iconIDs := make([]int, 0, len(icons))
	for id := range icons {
		if id != currentIconId {
			iconIDs = append(iconIDs, id)
		}
	}
	return iconIDs[rand.Intn(len(iconIDs))]
}

func getIconName(iconID int) string {
	return icons[iconID]
}

func sendVerificationMessage(e *handler.CommandEvent, accountData constants.AccountData, dataID string) error {
	imageURL := fmt.Sprintf("https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default//v1/profile-icons/%d.jpg", accountData.VerifyImageID)

	embed := discord.NewEmbedBuilder().
		SetTitle("Verify Account").
		SetDescription(fmt.Sprintf(
			"Account %s#%s has been found.\nChange your profile picture to the `%s` and click the verify button.",
			accountData.GameName, accountData.TagLine, getIconName(accountData.VerifyImageID)),
		).
		SetImage(imageURL).
		Build()

	actionRow := discord.NewPrimaryButton("Verify", "verify_acc/"+dataID)

	_, err := e.Client().Rest().UpdateInteractionResponse(
		e.ApplicationID(),
		e.Token(),
		discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			AddActionRow(actionRow).
			Build(),
	)

	if err != nil {
		return fmt.Errorf("failed to update interaction response: %w", err)
	}

	return nil
}
