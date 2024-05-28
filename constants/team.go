package constants

import (
	"github.com/bwmarrin/snowflake"
	dSnowflake "github.com/disgoorg/snowflake/v2"
)


type Match struct {
	MatchID   snowflake.ID   `json:"matchID"`
	Team1     []dSnowflake.ID `json:"team1"`
	Team2     []dSnowflake.ID `json:"team2"`
	Winner    string   `json:"winner"`
	Timestamp int64    `json:"timestamp"`
}
