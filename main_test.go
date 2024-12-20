package main

import (
	"encoding/json"
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestGetAllFailedUsers(t *testing.T) {
	failedUsersDB, err := leveldb.OpenFile("bots/bot5/nethealth.failedUsers", nil)
	if err != nil {
		panic(err)
	}
	iter := failedUsersDB.NewIterator(nil, nil)
	defer iter.Release()
	var users []*model.User
	for iter.Next() {
		u := model.User{}
		err := json.Unmarshal(iter.Value(), &u)
		if err != nil {
			t.Error(err)
		}
		users = append(users, &u)
	}
	if err := iter.Error(); err != nil {
		t.Error(err)
	}
	t.Logf("get %d failed users", len(users))
}
