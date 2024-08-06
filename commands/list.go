package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var list = discord.SlashCommandCreate{
	Name:        "list",
	Description: "Shows a list of all Profiles and their ratings.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "filter",
			Description: "Filter users by status (online or offline)",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "Online", Value: "online"},
				{Name: "Offline", Value: "offline"},
			},
		},
	},
}

func ListHandler(b *dbot.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		profiles, err := utils.GetProfiles()
		if err != nil {
			return err
		}

		guildID := e.GuildID()
		filter := e.SlashCommandInteractionData().String("filter")

		const maxEmbedLength = 2000
		const maxFieldsPerEmbed = 21
		var pages []discord.EmbedBuilder

		pageUsers := make([]constants.Profile, 0)
		currentLength := 0
		fieldCount := 0
		totalUsers := len(profiles)
		onlineUsers := 0
		displayedUsers := 0

		for _, profile := range profiles {
			isOnline := isUserOnline(b, *guildID, profile.UserID)
			if isOnline {
				onlineUsers++
			}

			// Apply filter if specified
			if filter == "online" && !isOnline {
				continue
			}
			if filter == "offline" && isOnline {
				continue
			}

			fieldValue := createFieldValue(profile, isOnline)

			if currentLength+len(fieldValue) > maxEmbedLength || fieldCount >= maxFieldsPerEmbed {
				embed := createEmbed(pageUsers, totalUsers, onlineUsers, displayedUsers, filter, b, *guildID)
				pages = append(pages, *embed)

				pageUsers = make([]constants.Profile, 0)
				currentLength = 0
				fieldCount = 0
			}

			pageUsers = append(pageUsers, profile)
			currentLength += len(fieldValue)
			fieldCount++
			displayedUsers++
		}

		if len(pageUsers) > 0 {
			embed := createEmbed(pageUsers, totalUsers, onlineUsers, displayedUsers, filter, b, *guildID)
			pages = append(pages, *embed)
		}

		return b.Paginator.Create(e.Respond, paginator.Pages{
			ID: e.ID().String(),
			PageFunc: func(page int, embed *discord.EmbedBuilder) {
				*embed = pages[page]
				embed.SetFooter(fmt.Sprintf("Page %d of %d", page+1, len(pages)), "")
			},
			Pages:      len(pages),
			ExpireMode: paginator.ExpireModeAfterLastUsage,
		}, false)
	}
}

func createEmbed(users []constants.Profile, totalUsers, onlineUsers, displayedUsers int, filter string, b *dbot.Bot, guildID snowflake.ID) *discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder().
		SetTitle("User Profiles").
		SetColor(0x3498db)

	var filterDescription string
	switch filter {
	case "online":
		filterDescription = fmt.Sprintf("Online: %d/%d (Showing %d online users)", onlineUsers, totalUsers, displayedUsers)
	case "offline":
		filterDescription = fmt.Sprintf("Offline: %d/%d (Showing %d offline users)", totalUsers-onlineUsers, totalUsers, displayedUsers)
	default:
		filterDescription = fmt.Sprintf("Online: %d/%d (Showing all %d users)", onlineUsers, totalUsers, displayedUsers)
	}
	embed.SetDescription(filterDescription)

	for _, profile := range users {
		isOnline := isUserOnline(b, guildID, profile.UserID)
		fieldValue := createFieldValue(profile, isOnline)
		embed.AddField("", fieldValue, true)
	}

	return embed
}

func createFieldValue(profile constants.Profile, isOnline bool) string {
	winrate := 0.0
	if total := profile.Wins + profile.Losses; total > 0 {
		winrate = float64(profile.Wins) / float64(total) * 100
	}

	status := ":red_circle:"
	if isOnline {
		status = ":green_circle:"
	}

	return fmt.Sprintf("<@%s> %s\nRating: %d\nW/L: %d/%d\nWinrate: %.2f%%",
		profile.UserID, status, profile.Rating, profile.Wins, profile.Losses, winrate)
}

func isUserOnline(b *dbot.Bot, guildID, userID snowflake.ID) bool {
	presence, ok := b.Client.Caches().Presence(guildID, userID)
	return ok && presence.Status != discord.OnlineStatusOffline
}
