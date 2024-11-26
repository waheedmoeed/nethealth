package model

type Transaction struct {
	ServiceDate string `json:"serviceDate"`
	ServiceCode string `json:"serviceCode"`
	Description string `json:"description"`
	ClaimType   string `json:"claimType"`
	Units       string `json:"units"`
	Rate        string `json:"rate"`
	Charge      string `json:"charge"`
	Payer       string `json:"payer"`
	Batch       string `json:"batch"`
	Balance     string `json:"balance"`
	Entity      string `json:"entity"`
}
