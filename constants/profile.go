package constants

import "github.com/disgoorg/snowflake/v2"

type Profile struct {
	UserID        snowflake.ID `json:"userID"`
	Rating        int          `json:"rating"`
	Wins          int          `json:"wins"`
	Losses        int          `json:"losses"`
	VerifiedPUUID *string      `json:"verifiedPUUID,omitempty"`
	SummonerID    *string      `json:"summonerID,omitempty"`
}
