package constants

type SummonerResponse struct {
	ProfileIconId int `json:"profileIconId"`
}

type RiotAccountResponse struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type AccountData struct {
	RiotAccountResponse
	Region        string
	VerifyImageID int
}
