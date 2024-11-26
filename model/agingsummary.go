package model

type AgingSummary struct {
	PayerPlan      string `json:"payerPlan"`
	Credits        string `json:"credits"`
	Total          string `json:"total"`
	Current        string `json:"current"`
	ThirtyDays     string `json:"thirtyDays"`
	SixtyDays      string `json:"sixtyDays"`
	NinetyDays     string `json:"ninetyDays"`
	OneTwentyDays  string `json:"oneTwentyDays"`
	OneEightyDays  string `json:"oneEightyDays"`
	MoreThanEigthy string `json:"moreThanEighty"`
}
