package leveldb

import (
	"encoding/json"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/syndtr/goleveldb/leveldb"
)

var levelDB *leveldb.DB

const path = "./nethealth.leveldb"

func init() {
	var err error
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	levelDB = db
}

func PutJob(job *model.Job) error {
	document, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return levelDB.Put([]byte(job.FileName), document, nil)
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
	return levelDB.Write(batch, nil)
}

func DeleteJob(fileName string) error {
	return levelDB.Delete([]byte(fileName), nil)
}

func GetJobs() ([]*model.Job, error) {
	iter := levelDB.NewIterator(nil, nil)
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
