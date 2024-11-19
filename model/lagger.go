package model

type LaggerGroup struct {
	Laggers              []*Lagger
	ServiceDate          string
	EstimatedAdjustments []*Lagger
}
type Lagger struct {
	TxDate        string `json:"txDate"`
	Type          string `json:"type"`
	ControlNumber string `json:"controlNumber"`
	Description   string `json:"description"`
	Seq           string `json:"seq"`
	ServiceDate   string `json:"serviceDate"`
	Category      string `json:"category"`
	DBAmount      string `json:"dbAmount"`
	CRAmount      string `json:"crAmount"`
	Balance       string `json:"balance"`
	PDFLink       string `json:"pdfLink"`
	UploadedLink  string `json:"uploadedLink"`
}
