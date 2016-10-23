package model

import (
	"os"
	"testing"
	"time"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestBuildBackupfileName(t *testing.T) {
	filename := BuildBackupfileName("/tmp", "mydb", time.Unix(1313123123, 0))
	if err := AssertThat(filename.String(), Is("/tmp/postgres_mydb_2011-08-12.dump")); err != nil {
		t.Fatal(err)
	}
}
