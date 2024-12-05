package leveldb

import (
	"testing"
)

func TestGetJobs(t *testing.T) {
	jobs, err := GetJobs()
	if err != nil {
		t.Error(err)
	}
	if len(jobs) == 0 {
		t.Error("no jobs found")
	}
	t.Logf("get %d jobs", 20)
}
