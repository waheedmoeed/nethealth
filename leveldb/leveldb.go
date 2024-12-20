package leveldb

import (
	"encoding/json"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var downloadJobs *leveldb.DB
var failedUsers *leveldb.DB
var latestAgencyState *leveldb.DB

var failedIter iterator.Iterator

const downloadJobPath = "./nethealth.downloadJobs"
const failedUserPath = "./nethealth.failedUsers"
const latestAgencyStatePath = "./nethealth.latestAgencyState"

func init() {
	var err error
	jobDB, err := leveldb.OpenFile(downloadJobPath, nil)
	if err != nil {
		panic(err)
	}
	failedUsersDB, err := leveldb.OpenFile(failedUserPath, nil)
	if err != nil {
		panic(err)
	}

	latestAgencyStateDB, err := leveldb.OpenFile(latestAgencyStatePath, nil)
	if err != nil {
		panic(err)
	}

	downloadJobs = jobDB
	failedUsers = failedUsersDB
	failedIter = failedUsers.NewIterator(nil, nil)
	latestAgencyState = latestAgencyStateDB
}

func PutJob(job *model.Job) error {
	document, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return downloadJobs.Put([]byte(job.FileName), document, nil)
}

func PutJobs(jobs []*model.Job) error {
	batch := new(leveldb.Batch)
	for _, job := range jobs {
		document, err := json.Marshal(job)
		if err != nil {
			return err
		}
		batch.Put([]byte(job.FileName), document)
	}
	return downloadJobs.Write(batch, nil)
}

func DeleteJob(fileName string) error {
	return downloadJobs.Delete([]byte(fileName), nil)
}

func GetJobs() ([]*model.Job, error) {
	iter := downloadJobs.NewIterator(nil, nil)
	defer iter.Release()
	var jobs []*model.Job
	for iter.Next() {
		if len(jobs) == 20 {
			break
		}
		j := model.Job{}
		err := json.Unmarshal(iter.Value(), &j)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &j)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func PutFailedUsers(users []*model.User) error {
	batch := new(leveldb.Batch)
	for _, user := range users {
		document, err := json.Marshal(user)
		if err != nil {
			return err
		}
		batch.Put([]byte(user.GetID()), document)
	}
	return failedUsers.Write(batch, nil)
}

func PutFailedUser(user *model.User) error {
	document, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return failedUsers.Put([]byte(user.GetID()), document, nil)
}

func HasFailedUsers() bool {
	iter := failedUsers.NewIterator(nil, nil)
	defer iter.Release()
	return iter.Next()
}

func GetFailedUsers() ([]*model.User, error) {
	var users []*model.User
	for failedIter.Next() {
		if len(users) == 5 {
			break
		}
		u := model.User{}
		err := json.Unmarshal(failedIter.Value(), &u)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := failedIter.Error(); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteFailedUser(user *model.User) error {
	return failedUsers.Delete([]byte(user.GetID()), nil)
}

func PutAgencyState(agency string, latestRecord string) error {
	return latestAgencyState.Put([]byte(agency), []byte(latestRecord), nil)
}

func GetAgencyState(agency string) (string, error) {
	value, err := latestAgencyState.Get([]byte(agency), nil)
	if err != nil {
		return "", err
	}
	return string(value), nil
}
