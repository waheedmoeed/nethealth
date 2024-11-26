package model

type TransactionDetail struct {
	ServiceDate          string                  `json:"serviceDate"`
	ServiceCode          string                  `json:"serviceCode"`
	Description          string                  `json:"description"`
	Charge               string                  `json:"charge"`
	Balance              string                  `json:"balance"`
	TransactionBreakdown []*TransactionBreakdown `json:"transactionBreakdown"`
}

type TransactionBreakdown struct {
	Date         string `json:"date"`
	ResonCode    string `json:"resonCode"`
	Description  string `json:"description"`
	Amount       string `json:"amount"`
	Reference    string `json:"reference"`
	Payer        string `json:"payer"`
	Batch        string `json:"batch"`
	PDFLink      string `json:"pdfLink"`
	UploadedLink string `json:"uploadedLink"`
}

func (t *TransactionDetail) GetDetailedName() string {
	return t.ServiceDate + " - " + t.ServiceCode + " - " + t.Description
}
