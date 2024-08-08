package constants

type SummonerResponse struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	ID            string `json:"id"`
	PUUID         string `json:"puuid"`
	SummonerLevel int    `json:"summonerLevel"`
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
