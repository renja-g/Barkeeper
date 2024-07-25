package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var list = discord.SlashCommandCreate{
	Name:        "list",
	Description: "Shows a list of all users and their ratings.",
}

func ListHandler(e *handler.CommandEvent, b *dbot.Bot) error {
    ratings, err := utils.GetRatings()
    if err != nil {
        return err
    }

    const maxEmbedLength = 2000
    const maxFieldsPerEmbed = 21
    var pages []*discord.EmbedBuilder
    
    pageUsers := make([]constants.Rating, 0)
    currentLength := 0
    fieldCount := 0
    
    for _, rating := range ratings {
        winrate := 0.0
        if rating.Wins+rating.Losses > 0 {
            winrate = float64(rating.Wins) / float64(rating.Wins+rating.Losses) * 100
        }
        fieldValue := fmt.Sprintf("<@%s>\nRating: %d\nW/L: %d/%d\nWinrate: %.2f%%", 
            rating.UserID, rating.Rating, rating.Wins, rating.Losses, winrate)
        
        if currentLength + len(fieldValue) > maxEmbedLength || fieldCount >= maxFieldsPerEmbed {
            embed := createEmbed(pageUsers)
            pages = append(pages, embed)
            
            pageUsers = make([]constants.Rating, 0)
            currentLength = 0
            fieldCount = 0
        }
        
        pageUsers = append(pageUsers, rating)
        currentLength += len(fieldValue)
        fieldCount++
    }
    
    if len(pageUsers) > 0 {
        embed := createEmbed(pageUsers)
        pages = append(pages, embed)
    }
    
    return b.Paginator.Create(e.Respond, paginator.Pages{
        ID: e.ID().String(),
        PageFunc: func(page int, embed *discord.EmbedBuilder) {
            *embed = *pages[page]
            embed.SetFooter(fmt.Sprintf("Page %d of %d", page+1, len(pages)), "")
        },
        Pages:      len(pages),
        ExpireMode: paginator.ExpireModeAfterLastUsage,
    }, false)
}

func createEmbed(users []constants.Rating) *discord.EmbedBuilder {
    embed := discord.NewEmbedBuilder().
        SetTitle("User Ratings").
        SetColor(0x3498db)

    for _, rating := range users {
        winrate := 0.0
        if rating.Wins+rating.Losses > 0 {
            winrate = float64(rating.Wins) / float64(rating.Wins+rating.Losses) * 100
        }
        fieldValue := fmt.Sprintf("<@%s>\nRating: %d\nW/L: %d/%d\nWinrate: %.2f%%", 
            rating.UserID, rating.Rating, rating.Wins, rating.Losses, winrate)
        embed.AddField("", fieldValue, true)
    }

    return embed
}
