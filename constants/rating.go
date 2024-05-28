package constants

import "github.com/disgoorg/snowflake/v2"

type Rating struct {
	UserID snowflake.ID `json:"userID"`
	Rating int    `json:"rating"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}
