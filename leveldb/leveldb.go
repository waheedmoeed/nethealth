package leveldb

import (
	"encoding/json"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/syndtr/goleveldb/leveldb"
)

var downloadJobs *leveldb.DB
var failedUsers *leveldb.DB

const downloadJobPath = "./nethealth.downloadJobs"
const failedUserPath = "./nethealth.failedUsers"

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

	downloadJobs = jobDB
	failedUsers = failedUsersDB
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

func PutFailedUser(user *model.User) error {
	document, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return failedUsers.Put([]byte(user.GetID()), document, nil)
}

func GetFailedUsers() ([]*model.User, error) {
	iter := failedUsers.NewIterator(nil, nil)
	defer iter.Release()
	var users []*model.User
	for iter.Next() {
		if len(users) == 20 {
			break
		}
		u := model.User{}
		err := json.Unmarshal(iter.Value(), &u)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return users, nil
}
