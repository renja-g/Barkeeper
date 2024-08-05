package constants

import (
	"github.com/disgoorg/snowflake/v2"
)

type Match struct {
	MatchID   snowflake.ID    `json:"matchID"`
	Team1     []snowflake.ID `json:"team1"`
	Team2     []snowflake.ID `json:"team2"`
	Winner    string          `json:"winner"`
	Timestamp int64           `json:"timestamp"`
}
