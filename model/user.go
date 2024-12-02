package model

import (
	"strconv"
	"strings"
)

const (
	BASE_DATA_DIR = "data"
	BUCKET_URL    = "https://storage.cloud.google.com/nethealth/"
)

type Users map[string]User
type User struct {
	ID               int    `json:"id"`
	UniqueIdentifier string `json:"uniqueIdentifier"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	MI               string `json:"mi"`
	AccountNumber    int64  `json:"accountNumber"`
	Enity            string `json:"entity"`
	IsMigrated       bool   `json:"isMigrated"`
}

func (user *User) GetID() string {
	return user.FirstName + "_" + user.LastName + "_" + strconv.FormatInt(user.AccountNumber, 10)
}
func (user *User) GetFullName() string {
	return user.FirstName + "_" + user.LastName
}

func (user *User) GetUserDataRoomPath() string {
	return BASE_DATA_DIR + "/" + strings.ReplaceAll(user.Enity, " ", "") + "/" + user.GetID()
}

func (user *User) GetJobFilePath() string {
	return BASE_DATA_DIR + "/jobs/" + user.GetID() + ".csv"
}

func (user *User) GetPendingJobFilePath() string {
	return BASE_DATA_DIR + "/jobs/" + user.GetID() + "_pending.csv"
}
