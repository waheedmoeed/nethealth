package keeper

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/abdulwaheed/nethealth/model"
)

// UsersDB stores users in memory and also writes them to a CSV file
type UsersDB struct {
	agencyName string
	users      map[string]*model.User
}

func NewUsersDB(agencyName string) (*UsersDB, error) {
	db := &UsersDB{
		agencyName: agencyName,
	}
	file, err := os.Open(db.agencyName + "_users.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for index, record := range records {
		accountNumber, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}
		isMigrated, err := strconv.ParseBool(record[4])
		if err != nil {
			return nil, err
		}
		db.users[record[0]] = &model.User{
			ID:            int64(index),
			FirstName:     record[1],
			LastName:      record[2],
			AccountNumber: accountNumber,
			IsMigrated:    isMigrated,
		}
	}

	return db, nil
}

func (db *UsersDB) GetUser(id string) (*model.User, error) {
	if user, ok := db.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (db *UsersDB) GetAllUsers() ([]*model.User, error) {
	users := make([]*model.User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}
	return users, nil
}
