package model

import (
	"strings"

	"github.com/abdulwaheed/nethealth/utils"
)

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

func (t *TransactionBreakdown) GetFileName() string {
	str := t.Date + "_" + t.ResonCode + "_" + t.Batch
	str = strings.ReplaceAll(str, "/", "_")
	str = utils.RemoveBetweenChars(str, "<", ">")
	return strings.ReplaceAll(str, " ", "")
}

func (t *TransactionBreakdown) GetFilePath(user *User) string {
	return user.GetUserDataRoomPath() + "/transactionbreakdowns/" + t.GetFileName()
}

func (t *TransactionBreakdown) GetUploadedLink(user *User) string {
	return BUCKET_URL + user.GetUserDataRoomPath() + "/transactionbreakdowns/" + t.GetFileName() + ".pdf"
}
