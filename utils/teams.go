package utils

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"sort"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/renja-g/Barkeeper/constants"
)

func GenerateTeams(users []constants.Profile) ([]constants.Profile, []constants.Profile) {
	n := len(users)
	if n%2 != 0 {
		return nil, nil
	}
	halfSize := n / 2

	var bestTeams [][2][]constants.Profile
	minDifference := math.MaxInt32

	// Helper function to calculate the difference in ratings between two teams
	calculateDifference := func(team1, team2 []constants.Profile) int {
		rating1, rating2 := 0, 0
		for _, profile := range team1 {
			rating1 += profile.Rating
		}
		for _, user := range team2 {
			rating2 += user.Rating
		}
		return abs(rating1 - rating2)
	}

	// Generate all possible team combinations using bitwise operations
	totalCombinations := 1 << n // 2^n
	for i := 0; i < totalCombinations; i++ {
		var team1, team2 []constants.Profile
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
				bestTeams = [][2][]constants.Profile{{team1, team2}}
			} else if difference == minDifference {
				bestTeams = append(bestTeams, [2][]constants.Profile{team1, team2})
			}
		}
	}

	// Check if there are any best teams found
	if len(bestTeams) == 0 {
		return nil, nil
	}

	// Randomly select one of the best teams
	selectedIndex := rand.Intn(len(bestTeams))
	bestTeam1, bestTeam2 := bestTeams[selectedIndex][0], bestTeams[selectedIndex][1]

	return bestTeam1, bestTeam2
}

func FormatTeam(team []constants.Profile) string {
	var formatted string
	for i, user := range team {
		if i != 0 {
			formatted += "\n"
		}
		formatted += fmt.Sprintf("<@%s> %d", user.UserID, user.Rating)
	}
	return formatted
}

func CalculateTeamRating(team []constants.Profile) int {
	rating := 0
	for _, user := range team {
		rating += user.Rating
	}
	return rating
}

func ParseTeamMessage(message discord.Message) ([]snowflake.ID, []snowflake.ID) {
	var team1, team2 []snowflake.ID
	re := regexp.MustCompile(`<@(\d+)> \d+`)
	team1Matches := re.FindAllStringSubmatch(message.Embeds[0].Fields[0].Value, -1)
	team2Matches := re.FindAllStringSubmatch(message.Embeds[0].Fields[1].Value, -1)

	for _, match := range team1Matches {
		if len(match) < 2 {
			log.Printf("Unexpected match format for team1: %v", match)
			continue
		}
		id, err := snowflake.Parse(match[1])
		if err != nil {
			log.Printf("Error parsing ID for team1: %v", err)
			continue
		}
		team1 = append(team1, id)
	}

	for _, match := range team2Matches {
		if len(match) < 2 {
			log.Printf("Unexpected match format for team2: %v", match)
			continue
		}
		id, err := snowflake.Parse(match[1])
		if err != nil {
			log.Printf("Error parsing ID for team2: %v", err)
			continue
		}
		team2 = append(team2, id)
	}
	log.Printf("team1: %v", team1)
	log.Printf("team2: %v", team2)
	return team1, team2
}


func ParseMatchID(message discord.Message) (snowflake.ID, error) {
	matchID := message.Embeds[0].Footer.Text[9:28]
	return snowflake.Parse(matchID)
}

func SortProfilesByWinRate(profiles []constants.Profile) []constants.Profile {
	sort.Slice(profiles, func(i, j int) bool {
		winRate1, winRate2 := 0.0, 0.0
		if profiles[i].Wins+profiles[i].Losses != 0 {
			winRate1 = float64(profiles[i].Wins) / float64(profiles[i].Wins+profiles[i].Losses)
		}
		if profiles[j].Wins+profiles[j].Losses != 0 {
			winRate2 = float64(profiles[j].Wins) / float64(profiles[j].Wins+profiles[j].Losses)
		}
		return winRate1 > winRate2
	})

	return profiles
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
