package ufcstats

type JudgeScore struct {
	Name      string `json:"name"`
	RedScore  int    `json:"red_score"`
	BlueScore int    `json:"blue_score"`
}

type FightBonus struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient,omitempty"`
}

type RoundFighterStat struct {
	Round int `json:"round"`

	KD             int    `json:"kd"`
	SigLanded      int    `json:"sig_landed"`
	SigAttempted   int    `json:"sig_attempted"`
	TotalLanded    int    `json:"total_landed"`
	TotalAttempted int    `json:"total_attempted"`
	TDLanded       int    `json:"td_landed"`
	TDAttempted    int    `json:"td_attempted"`
	SubAtt         int    `json:"sub_att"`
	Rev            int    `json:"rev"`
	CTRL           string `json:"ctrl"`

	HeadLanded        int `json:"head_landed"`
	HeadAttempted     int `json:"head_attempted"`
	BodyLanded        int `json:"body_landed"`
	BodyAttempted     int `json:"body_attempted"`
	LegLanded         int `json:"leg_landed"`
	LegAttempted      int `json:"leg_attempted"`
	DistanceLanded    int `json:"distance_landed"`
	DistanceAttempted int `json:"distance_attempted"`
	ClinchLanded      int `json:"clinch_landed"`
	ClinchAttempted   int `json:"clinch_attempted"`
	GroundLanded      int `json:"ground_landed"`
	GroundAttempted   int `json:"ground_attempted"`
}

type FightRoundStats struct {
	Red  []RoundFighterStat `json:"red"`
	Blue []RoundFighterStat `json:"blue"`
}

type FightDetailsScrape struct {
	Fight       *Fight          `json:"fight"`
	Red         *Fighter        `json:"red,omitempty"`
	Blue        *Fighter        `json:"blue,omitempty"`
	IsTitleBout bool            `json:"is_title_bout"`
	RefereeName string          `json:"referee_name"`
	Bonuses     []FightBonus    `json:"bonuses"`
	Judges      []JudgeScore    `json:"judges"`
	RoundStats  FightRoundStats `json:"round_stats"`
	Rounds      int             `json:"rounds"`
}
