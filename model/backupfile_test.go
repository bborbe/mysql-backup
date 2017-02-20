package model

import (
	"os"
	"testing"
	"time"

	"io/ioutil"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestBuildBackupfileName(t *testing.T) {
	filename := BuildBackupfileName("myname", "/tmp", "mydb", time.Unix(1313123123, 0))
	if err := AssertThat(filename.String(), Is("/tmp/myname_mydb_2011-08-12.dump")); err != nil {
		t.Fatal(err)
	}
}

func TestExistsReturnTrueIfExistsAndNotEmpty(t *testing.T) {
	file, err := ioutil.TempFile("", "backupfile")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("hello world")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	file.Close()

	b := BackupFilename(file.Name())
	if err := AssertThat(b.Exists(), Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestExistsReturnFalseIfExistsButEmpty(t *testing.T) {
	file, err := ioutil.TempFile("", "backupfile")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Close()

	b := BackupFilename(file.Name())
	if err := AssertThat(b.Exists(), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestExistsReturnFalseIfNotExisting(t *testing.T) {
	b := BackupFilename("/filedoesnotexists")
	if err := AssertThat(b.Exists(), Is(false)); err != nil {
		t.Fatal(err)
	}
}
