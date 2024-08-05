package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/renja-g/Barkeeper/constants"
	"github.com/renja-g/Barkeeper/utils"
)

var rate = discord.SlashCommandCreate{
	Name:        "rate",
	Description: "Rates a user.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to rate.",
			Required:    true,
		},
		discord.ApplicationCommandOptionInt{
			Name:         "rating",
			Description:  "The rating to give the user.",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func RateHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		embed := discord.NewEmbedBuilder().
			SetTitle("Rating set").
			SetDescriptionf("Rating for %s set to %d", e.SlashCommandInteractionData().User("user").Mention(), e.SlashCommandInteractionData().Int("rating")).
			SetColor(0x00ff00).
			Build()

		ratings, err := utils.GetRatings()
		if err != nil {
			return err
		}

		// Check if the user has already been rated
		found := false
		for i, rating := range ratings {
			if rating.UserID == e.SlashCommandInteractionData().User("user").ID {
				ratings[i].Rating = e.SlashCommandInteractionData().Int("rating")
				found = true
				break
			}
		}

		// If the user has not been rated yet, add a new rating
		if !found {
			ratings = append(ratings, constants.Rating{
				UserID: e.SlashCommandInteractionData().User("user").ID,
				Rating: e.SlashCommandInteractionData().Int("rating"),
				// Wins and losses are set to 0 by default
			})
		}

		err = utils.SaveRatings(ratings)
		if err != nil {
			return err
		}

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}
}

func RateAutocompleteHandler(e *handler.AutocompleteEvent) error {
	return e.AutocompleteResult([]discord.AutocompleteChoice{
		discord.AutocompleteChoiceInt{
			Name:  "1",
			Value: 1,
		},
		discord.AutocompleteChoiceInt{
			Name:  "2",
			Value: 2,
		},
		discord.AutocompleteChoiceInt{
			Name:  "3",
			Value: 3,
		},
		discord.AutocompleteChoiceInt{
			Name:  "4",
			Value: 4,
		},
		discord.AutocompleteChoiceInt{
			Name:  "5",
			Value: 5,
		},
		discord.AutocompleteChoiceInt{
			Name:  "6",
			Value: 6,
		},
		discord.AutocompleteChoiceInt{
			Name:  "7",
			Value: 7,
		},
		discord.AutocompleteChoiceInt{
			Name:  "8",
			Value: 8,
		},
		discord.AutocompleteChoiceInt{
			Name:  "9",
			Value: 9,
		},
		discord.AutocompleteChoiceInt{
			Name:  "10",
			Value: 10,
		},
	})
}
