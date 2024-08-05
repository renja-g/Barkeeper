package utils

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"

	"github.com/bwmarrin/snowflake"
	"github.com/disgoorg/disgo/discord"
	dSnowflake "github.com/disgoorg/snowflake/v2"

	"github.com/renja-g/Barkeeper/constants"
)

func GenerateTeams(users []constants.Rating) ([]constants.Rating, []constants.Rating) {
	n := len(users)
	if n%2 != 0 {
		return nil, nil
	}
	halfSize := n / 2

	var bestTeams [][2][]constants.Rating
	minDifference := math.MaxInt32

	// Helper function to calculate the difference in ratings between two teams
	calculateDifference := func(team1, team2 []constants.Rating) int {
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
		var team1, team2 []constants.Rating
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
				bestTeams = [][2][]constants.Rating{{team1, team2}}
			} else if difference == minDifference {
				bestTeams = append(bestTeams, [2][]constants.Rating{team1, team2})
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

func FormatTeam(team []constants.Rating) string {
	var formatted string
	for i, user := range team {
		if i != 0 {
			formatted += "\n"
		}
		formatted += fmt.Sprintf("<@%s> %d", user.UserID, user.Rating)
	}
	return formatted
}

func CalculateTeamRating(team []constants.Rating) int {
	rating := 0
	for _, user := range team {
		rating += user.Rating
	}
	return rating
}

func ParseTeamMessage(message discord.Message) ([]dSnowflake.ID, []dSnowflake.ID) {
	var team1, team2 []dSnowflake.ID
	re := regexp.MustCompile(`<@(\d+)> \d+`)
	team1Matches := re.FindAllString(message.Embeds[0].Fields[0].Value, -1)
	team2Matches := re.FindAllString(message.Embeds[0].Fields[1].Value, -1)

	for _, match := range team1Matches {
		id, err := dSnowflake.Parse(match[2:20])
		if err != nil {
			continue
		}
		team1 = append(team1, id)
	}

	for _, match := range team2Matches {
		id, err := dSnowflake.Parse(match[2:20])
		if err != nil {
			continue
		}
		team2 = append(team2, id)
	}

	return team1, team2
}

func ParseMatchID(message discord.Message) (snowflake.ID, error) {
	matchID := message.Embeds[0].Footer.Text[9:28]
	return snowflake.ParseString(matchID)
}

func SortRatingsByWinRate(ratings []constants.Rating) []constants.Rating {
	sort.Slice(ratings, func(i, j int) bool {
		winRate1, winRate2 := 0.0, 0.0
		if ratings[i].Wins+ratings[i].Losses != 0 {
			winRate1 = float64(ratings[i].Wins) / float64(ratings[i].Wins+ratings[i].Losses)
		}
		if ratings[j].Wins+ratings[j].Losses != 0 {
			winRate2 = float64(ratings[j].Wins) / float64(ratings[j].Wins+ratings[j].Losses)
		}
		return winRate1 > winRate2
	})

	return ratings
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
