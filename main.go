package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/joho/godotenv"
)

// Global variables
var (
    GuildID        string
    BotToken       string
    RemoveCommands = true
	Node           *snowflake.Node
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	godotenv.Load(".env")
	GuildID = os.Getenv("GUILD_ID")
    BotToken = os.Getenv("BOT_TOKEN")

	var err error

	// Create a new Node with a Node number of 1
	Node, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue = 0.0

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "rate",
			Description: "Rates of sets the rating of a user",
			// The rating will be saved in a json file ratings.json
			// If the user is already in the file, the rating will be updated
			/*
				[
					{
						"userID": 1234,
						"rating": 0,
						"wins": 0,
						"loses": 0
					}
				]
			*/
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to rate",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "rating",
					Description: "Rating",
					Required:    true,
					MaxValue:    10,
					MinValue:    &integerOptionMinValue,
				},
			},
		},
		{
			Name:        "info",
			Description: "Returns information about a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to get information about",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "Lists all the rated users",
		},
		{
			Name:        "teams",
			Description: "Generates fair teams of the users in the voice channel",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"rate": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Get the user ID and rating from the options
			userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
			rating := i.ApplicationCommandData().Options[1].IntValue()

			// Load the ratings from the file
			ratings, err := loadRatings()
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot load ratings",
				})
				return
			}

			// Find the user in the ratings
			var userRating *ratingData
			for _, r := range ratings {
				if r.UserID == userID {
					userRating = r
					break
				}
			}

			// If the user is not in the ratings, create a new entry
			if userRating == nil {
				userRating = &ratingData{
					// When a new ratingData struct is created, Wins and Losses are initialized to 0
					UserID: userID,
				}
				ratings = append(ratings, userRating)
			}

			// Update the rating
			userRating.Rating = int(rating)

			// Save the ratings back to the file
			err = saveRatings(ratings)
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot save ratings",
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Rating set",
							Description: fmt.Sprintf("Rating for <@%s> set to %d", userID, rating),
							Color:       0x00ff00,
						},
					},
				},
			})
		},
		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Get the user ID from the options
			userID := i.ApplicationCommandData().Options[0].UserValue(s).ID

			// Load the ratings from the file
			ratings, err := loadRatings()
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot load ratings",
				})
				return
			}

			// Find the user in the ratings
			var userRating *ratingData
			for _, r := range ratings {
				if r.UserID == userID {
					userRating = r
					break
				}
			}

			// If the user is not in the ratings return an error
			if userRating == nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Error",
								Description: "User not found",
								Color:       0xff0000,
							},
						},
					},
				})

				return
			}

			winrate := 0.0
			if userRating.Wins+userRating.Loses > 0 {
				winrate = float64(userRating.Wins) / float64(userRating.Wins+userRating.Loses) * 100
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "User information",
							Description: fmt.Sprintf("User: <@%s>\nRating: %d\nWins: %d\nLoses: %d\nWinrate: %.2f%%", userID, userRating.Rating, userRating.Wins, userRating.Loses, winrate),
							Color:       0x00ff00,
						},
					},
				},
			})
		},
		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Load the ratings from the file
			ratings, err := loadRatings()
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot load ratings",
				})
				return
			}
		
			// Calculate winrate and store it in the struct
			for _, r := range ratings {
				if r.Wins+r.Loses > 0 {
					r.Winrate = float64(r.Wins) / float64(r.Wins+r.Loses) * 100
				} else {
					r.Winrate = 0.0
				}
			}
		
			// Sort the ratings by winrate
			sort.Slice(ratings, func(i, j int) bool {
				return ratings[i].Winrate > ratings[j].Winrate
			})
		
			// Generate the list of users
			fields := []*discordgo.MessageEmbedField{}
			for _, r := range ratings {
				userField := &discordgo.MessageEmbedField{
					Value:  fmt.Sprintf("<@%s>\nRating: %d\nWins: %d\nLosses: %d\nWinrate: %.2f%%", r.UserID, r.Rating, r.Wins, r.Loses, r.Winrate),
					Inline: true,
				}
				fields = append(fields, userField)
			}
		
			// Send the list
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:  "Ratings",
							Fields: fields,
							Color:  0x00ff00,
						},
					},
				},
			})
		},
		"teams": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, err := s.State.Guild(i.GuildID)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Cannot get the guild",
					},
				})
				return
			}

			// Get the voice channel of the user who triggered the command
			var ChannelID string
			for _, vs := range guild.VoiceStates {
				if vs.UserID == i.Member.User.ID {
					ChannelID = vs.ChannelID
					break
				}
			}

			if ChannelID == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You must be in a voice channel to use this command",
					},
				})
				return
			}

			// Get the userIds in the voice channel
			var ids []string
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == ChannelID {
					ids = append(ids, vs.UserID)
				}
			}

			// Check if there are exactly 10 users in the voice channel
			/*
				if len(ids) != 10 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "There must be exactly 10 users in the voice channel",
						},
					})
					return
				}
			*/

			// Check if all users are in the ratings
			ratings, err := loadRatings()
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error loading ratings",
					},
				})
				return
			}

			var missingRatings []string
			for _, id := range ids {
				found := false
				for _, r := range ratings {
					if r.UserID == id {
						found = true
						break
					}
				}
				if !found {
					missingRatings = append(missingRatings, id)
				}
			}

			if len(missingRatings) > 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("The following users are missing ratings: %v", missingRatings),
					},
				})
				return
			}

			// drop all ratings that are not in the voice channel
			var currentUsers []*ratingData
			for _, r := range ratings {
				for _, id := range ids {
					if r.UserID == id {
						currentUsers = append(currentUsers, r)
						break
					}
				}
			}

			// Generate the teams
			team1, team2 := generateTeams(currentUsers)

			// Send the teams
			team1String := ""
			team1Rating := 0
			for _, r := range team1 {
				team1String += fmt.Sprintf("<@%s> %d\n", r.UserID, r.Rating)
				team1Rating += r.Rating
			}

			team2String := ""
			team2Rating := 0
			for _, r := range team2 {
				team2String += fmt.Sprintf("<@%s> %d\n", r.UserID, r.Rating)
				team2Rating += r.Rating
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Teams",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  fmt.Sprintf("Team 1 (%d)", team1Rating),
									Value: team1String,
								},
								{
									Name:  fmt.Sprintf("Team 2 (%d)", team1Rating),
									Value: team2String,
								},
							},
							Color: 0x00ff00,
						},
					},
					Flags: discordgo.MessageFlagsLoading,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "start",
									Emoji: &discordgo.ComponentEmoji{
										Name: "▶️",
									},
									Label: "Start Match",
									Style: discordgo.SecondaryButton,
								},
								discordgo.Button{
									CustomID: "reshuffle",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🔁",
									},
									Label: "Reshuffle Teams",
									Style: discordgo.SecondaryButton,
								},
							},
						},
					},
				},
			})
		},
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"reshuffle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Load the ratings from the file
			ratings, err := loadRatings()
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot load ratings",
				})
				return
			}

			// Get the voice channel of the user who triggered the command
			guild, err := s.State.Guild(i.GuildID)
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot get the guild",
				})
				return
			}

			var ChannelID string
			for _, vs := range guild.VoiceStates {
				if vs.UserID == i.Member.User.ID {
					ChannelID = vs.ChannelID
					break
				}
			}

			if ChannelID == "" {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "You must be in a voice channel to use this command",
				})
				return
			}

			// Get the userIds in the voice channel
			var ids []string
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == ChannelID {
					ids = append(ids, vs.UserID)
				}
			}

			// drop all ratings that are not in the voice channel
			var currentUsers []*ratingData
			for _, r := range ratings {
				for _, id := range ids {
					if r.UserID == id {
						currentUsers = append(currentUsers, r)
						break
					}
				}
			}

			// Generate the teams
			team1, team2 := generateTeams(currentUsers)

			team1String := ""
			team1Rating := 0
			for _, r := range team1 {
				team1String += fmt.Sprintf("<@%s> %d\n", r.UserID, r.Rating)
				team1Rating += r.Rating
			}

			team2String := ""
			team2Rating := 0
			for _, r := range team2 {
				team2String += fmt.Sprintf("<@%s> %d\n", r.UserID, r.Rating)
				team2Rating += r.Rating
			}

			// Edit the message
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Teams",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  fmt.Sprintf("Team 1 (%d)", team1Rating),
									Value: team1String,
								},
								{
									Name:  fmt.Sprintf("Team 2 (%d)", team1Rating),
									Value: team2String,
								},
							},
							Color: 0x00ff00,
							Footer: &discordgo.MessageEmbedFooter{
								Text: "Teams reshuffled",
							},
						},
					},
					Flags: discordgo.MessageFlagsLoading,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "start",
									Emoji: &discordgo.ComponentEmoji{
										Name: "▶️",
									},
									Label: "Start Match",
									Style: discordgo.SecondaryButton,
								},
								discordgo.Button{
									CustomID: "reshuffle",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🔁",
									},
									Label: "Reshuffle Teams",
									Style: discordgo.SecondaryButton,
								},
							},
						},
					},
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
			}
		},
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// save the match to the matches.json file
			// [
			// 	{
			// 		"team1": ["userID1", "userID2", "userID3", "userID4", "userID5"],
			// 		"team2": ["userID6", "userID7", "userID8", "userID9", "userID10"],
			// 		"winner": "",
			// 		"timestamp": 1234567890
			// 	}

			// Load the matches from the file
			matches, err := loadMatches()
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Cannot load matches",
				})
				return
			}

			// use the message to get the team members
			// parse the embed fields to get the team members
			team1, team2 := parseTeams(i)

			match := &matchData{
				MatchID: Node.Generate().String(),
				Team1: team1,
				Team2: team2,
				Winner: "",
				Timestamp: time.Now().Unix(),
			}

			matches = append(matches, match)
			saveMatches(matches)

			// Edit the message
			oldEmbed := i.Message.Embeds[0]
			oldEmbed.Title = "Match in progress"
			oldEmbed.Description = "Select the winner of the match"
			oldEmbed.Footer = &discordgo.MessageEmbedFooter{
				Text: "MatchID: " + match.MatchID,
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						oldEmbed,
					},
					Flags: discordgo.MessageFlagsLoading,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "team1_wins",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🏆",
									},
									Label: "Team 1 wins",
									Style: discordgo.SuccessButton,
								},
								discordgo.Button{
									CustomID: "team2_wins",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🏆",
									},
									Label: "Team 2 wins",
									Style: discordgo.SuccessButton,
								},
								discordgo.Button{
									CustomID: "cancel_match",
									Emoji: &discordgo.ComponentEmoji{
										Name: "❌",
									},
									Label: "Cancel match",
									Style: discordgo.DangerButton,
								},
							},
						},
					},
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
			}
		},
		"team1_wins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},
		"team2_wins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},
		"cancel_match": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}

	})
}

type ratingData struct {
	UserID string `json:"userID"`
	Rating int    `json:"rating"`
	Wins   int    `json:"wins"`
	Loses  int    `json:"loses"`
	Winrate float64 // no json tag as this is calculated, not loaded
}

type matchData struct {
	MatchID   string   `json:"matchID"` // snowflake
	Team1     []string `json:"team1"`
	Team2     []string `json:"team2"`
	Winner    string   `json:"winner"`
	Timestamp int64    `json:"timestamp"`
}

func loadRatings() ([]*ratingData, error) {
	// Open the file
	file, err := os.Open("ratings.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a slice of ratingData
	var ratings []*ratingData
	err = json.Unmarshal(bytes, &ratings)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func loadMatches() ([]*matchData, error) {
	// Open the file
	file, err := os.Open("matches.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a slice of matchData
	var matches []*matchData
	err = json.Unmarshal(bytes, &matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func saveRatings(ratings []*ratingData) error {
	// Marshal the ratings into JSON
	bytes, err := json.MarshalIndent(ratings, "", "    ")
	if err != nil {
		return err
	}

	// Write the JSON to the file
	err = os.WriteFile("ratings.json", bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func saveMatches(matches []*matchData) error {
	// Marshal the matches into JSON
	bytes, err := json.MarshalIndent(matches, "", "    ")
	if err != nil {
		return err
	}

	// Write the JSON to the file
	err = os.WriteFile("matches.json", bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseTeams(i *discordgo.InteractionCreate) ([]string, []string) {
	// Define a regex to match user IDs in the format <@123456789> and optionally capture the score
	re := regexp.MustCompile(`<@(\d+)> \d+`)

	// Parse the embed fields to get the team members
	team1Matches := re.FindAllStringSubmatch(i.Message.Embeds[0].Fields[0].Value, -1)
	team2Matches := re.FindAllStringSubmatch(i.Message.Embeds[0].Fields[1].Value, -1)

	// Extract the user IDs
	extractUserIDs := func(matches [][]string) []string {
		userIDs := make([]string, len(matches))
		for i, match := range matches {
			userIDs[i] = match[1]
		}
		return userIDs
	}

	team1IDs := extractUserIDs(team1Matches)
	team2IDs := extractUserIDs(team2Matches)

	return team1IDs, team2IDs
}

func generateTeams(users []*ratingData) ([]*ratingData, []*ratingData) {
	n := len(users)
	halfSize := n / 2

	var bestTeams [][2][]*ratingData
	minDifference := math.MaxInt32

	// Helper function to calculate the difference in ratings between two teams
	calculateDifference := func(team1, team2 []*ratingData) int {
		rating1, rating2 := 0, 0
		for _, user := range team1 {
			rating1 += user.Rating
		}
		for _, user := range team2 {
			rating2 += user.Rating
		}
		return abs(rating1 - rating2)
	}

	// Generate all possible team combinations using bitwise operations
	totalCombinations := 1 << n // 2^n
	for i := 0; i < totalCombinations; i++ {
		var team1, team2 []*ratingData
		for j := 0; j < n; j++ {
			if i&(1<<j) != 0 {
				team1 = append(team1, users[j])
			} else {
				team2 = append(team2, users[j])
			}
		}
		if len(team1) == halfSize && len(team2) == halfSize {
			difference := calculateDifference(team1, team2)
			if difference < minDifference {
				minDifference = difference
				bestTeams = [][2][]*ratingData{{team1, team2}}
			} else if difference == minDifference {
				bestTeams = append(bestTeams, [2][]*ratingData{team1, team2})
			}
		}
	}

	// Randomly select one of the best teams
	selectedIndex := rand.Intn(len(bestTeams))
	bestTeam1, bestTeam2 := bestTeams[selectedIndex][0], bestTeams[selectedIndex][1]

	return bestTeam1, bestTeam2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if RemoveCommands {
		log.Println("Removing commands...")
		// We need to fetch the commands, since deleting requires the command ID.
		// We are doing this from the returned commands on line 375, because using
		// this will delete all the commands, which might not be desirable, so we
		// are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
