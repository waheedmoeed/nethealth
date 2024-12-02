package model

import (
	"strings"
)

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

func (l *Lagger) GetFileName() string {
	return strings.ReplaceAll(l.TxDate+"_"+l.Type, "/", "_")
}

func (l *Lagger) GetFilePath(user *User) string {
	return user.GetUserDataRoomPath() + "/laggers/" + l.GetFileName()
}

func (l *Lagger) GetUploadedLink(user *User) string {
	return BUCKET_URL + user.GetUserDataRoomPath() + "/laggers/" + l.GetFileName() + ".pdf"
}


func (l *Lagger) GetAdjustmentFileName() string {
	return strings.ReplaceAll(l.TxDate+"_"+l.Type+"_"+l.ControlNumber, "/", "_")
}

func (l *Lagger) GetAdjustmentFilePath(user *User) string {
	return user.GetUserDataRoomPath() + "/laggers/" + l.GetAdjustmentFileName()
}

func (l *Lagger) GetAdjustmentUploadedLink(user *User) string {
	return BUCKET_URL + user.GetUserDataRoomPath() + "/laggers/" + l.GetAdjustmentFileName() + ".pdf"
}