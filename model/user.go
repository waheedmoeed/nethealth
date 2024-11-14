package model

import "strconv"

type Users map[string]User
type User struct {
	ID               int `json:"id"`
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
