package model

import (
	"strings"

	"github.com/abdulwaheed/nethealth/utils"
)

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

func (c *Claim) GetFileName() string {
	str := c.ClaimNumber + "_" + c.PayingAgency
	str = strings.ReplaceAll(str, "/", "_")
	str = utils.RemoveBetweenChars(str, "<", ">")
	return strings.ReplaceAll(str, " ", "")
}

func (c *Claim) GetFilePath(user *User) string {
	return user.GetUserDataRoomPath() + "/claims/" + c.GetFileName()
}

func (c *Claim) GetUploadedLink(user *User) string {
	return BUCKET_URL + user.GetUserDataRoomPath() + "/claims/" + c.GetFileName() + ".pdf"
}
