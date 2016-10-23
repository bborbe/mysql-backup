package backup

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
	targetDirectory := "/tmp"
	database := "mydb"
	date := time.Unix(1313123123, 0)
	filename := buildBackupfileName(targetDirectory, database, date)
	if err := AssertThat(filename, Is("/tmp/postgres_mydb_2011-08-12.dump")); err != nil {
		t.Fatal(err)
	}
}

func TestRunCommand(t *testing.T) {
	err := runCommand("ls", "/")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
