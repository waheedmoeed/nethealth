package model

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	BASE_DATA_DIR = "data"
	BUCKET_URL    = "https://storage.cloud.google.com/nethealth/"
)

type Users map[string]User
type User struct {
	ID               int64  `json:"id"`
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

func (user *User) GetLedgerPageURL() string {
	return fmt.Sprintf("https://p13006.therapy.nethealth.com/Financials#patient/details/%d/ledger", user.AccountNumber)
}

func (user *User) GetClaimsPageURL() string {
	return fmt.Sprintf("https://p13006.therapy.nethealth.com/Financials#patient/details/%d/claims", user.AccountNumber)
}

func (user *User) GetTransactionsPageURL() string {
	return fmt.Sprintf("https://p13006.therapy.nethealth.com/Financials#patient/details/%d/transactions", user.AccountNumber)
}

func (user *User) GetBenefitsPageURL() string {
	return fmt.Sprintf("https://p13006.therapy.nethealth.com/Financials#patient/details/%d/benefitsVerification", user.AccountNumber)
}

func (user *User) GetAgingSummaryPageURL() string {
	return fmt.Sprintf("https://p13006.therapy.nethealth.com/Financials#patient/details/%d/agingSummary", user.AccountNumber)
}

func newUserFromCSVRecord(record []string, fileName string) (*User, error) {
	if len(record) != 7 { // assuming 6 fields in the CSV record
		return nil, errors.New("invalid CSV record")
	}

	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return nil, err
	}

	accountNumber, err := strconv.ParseInt(record[5], 10, 64)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:            id,
		FirstName:     record[2],
		LastName:      record[3],
		MI:            record[4],
		AccountNumber: accountNumber,
		IsMigrated:    record[5] == "true",
		Enity:         fileName,
	}

	return user, nil
}

func ReadUsersFromCSVFile(ctx context.Context, filePath string, fileName string) ([]*User, error) {
	users := make([]*User, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record[0] == "FirstName" {
			continue
		}
		user, err := newUserFromCSVRecord(record, fileName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}