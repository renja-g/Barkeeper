package constants

type Rating struct {
	UserID string `json:"userID"`
	Rating int    `json:"rating"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}
