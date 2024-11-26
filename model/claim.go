package model

type Claim struct {
	CreationDate    string `json:"creationDate"`
	ServicesFrom    string `json:"servicesFrom"`
	ServicesThrough string `json:"servicesThrough"`
	ClaimNumber     string `json:"claimNumber"`
	ClaimType       string `json:"claimType"`
	BatchNumber     string `json:"batchNumber"`
	Entity          string `json:"entity"`
	PayingAgency    string `json:"payingAgency"`
	PayerPlan       string `json:"payerPlan"`
	PayerSequence   string `json:"payerSequence"`
	ClaimAmount     string `json:"claimAmount"`
	PDFLink         string `json:"pdfLink"`
	UploadedLink    string `json:"uploadedLink"`
}
