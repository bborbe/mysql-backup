package backup

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/mysql_backup_cron/model"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func createBackuper() *Backuper {
	return New(
		"database",
		false,
		"localhost",
		3306,
		"user",
		"password",
		"name",
		"/backup",
		true,
		false,
	)
}
func createBackupAll() *Backuper {
	return New(
		"",
		true,
		"localhost",
		3306,
		"user",
		"password",
		"name",
		"/backup",
		true,
		false,
	)
}

func TestCreateDumper(t *testing.T) {
	b := createBackuper()
	if err := AssertThat(b.Database.String(), Is("database")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.AllDatabases, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Host.String(), Is("localhost")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Port.Int(), Is(3306)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.User.String(), Is("user")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Password.String(), Is("password")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Name.String(), Is("name")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.TargetDirectory.String(), Is("/backup")); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDumperAll(t *testing.T) {
	b := createBackupAll()
	if err := AssertThat(b.Database.String(), Is("")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.AllDatabases, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Host.String(), Is("localhost")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Port.Int(), Is(3306)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.User.String(), Is("user")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Password.String(), Is("password")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.Name.String(), Is("name")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(b.TargetDirectory.String(), Is("/backup")); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccess(t *testing.T) {
	b := createBackuper()
	if err := AssertThat(b.Validate(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccessAll(t *testing.T) {
	b := createBackupAll()
	if err := AssertThat(b.Validate(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateReturnErrorIfDatabaseIsEmptyAndAllDatabasesActiv(t *testing.T) {
	b := createBackupAll()
	b.Database = "mydb"
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfDatabaseIsEmpty(t *testing.T) {
	b := createBackuper()
	b.Database = ""
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfHostIsEmpty(t *testing.T) {
	b := createBackuper()
	b.Host = ""
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfPortIsZero(t *testing.T) {
	b := createBackuper()
	b.Port = 0
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfUserIsEmpty(t *testing.T) {
	mysqlDumper := createBackuper()
	mysqlDumper.User = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfPasswordIsEmpty(t *testing.T) {
	b := createBackuper()
	b.Password = ""
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfNameIsEmpty(t *testing.T) {
	b := createBackuper()
	b.Name = ""
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfTargetDirectoryIsEmpty(t *testing.T) {
	b := createBackuper()
	b.TargetDirectory = ""
	if err := AssertThat(b.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupShouldBeSkippedReturnAlwaysFalseIfOverwriteTrue(t *testing.T) {
	b := createBackuper()
	b.OverwriteBackup = true
	result := b.checkBackupShouldBeSkipped(model.BackupFilename("/tmp/backup"))
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupShouldBeSkippedReturnFalseIfFileNotExists(t *testing.T) {
	b := createBackuper()
	b.OverwriteBackup = false
	file, err := ioutil.TempFile("", "backupfile_not_exists")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	result := b.checkBackupShouldBeSkipped(model.BackupFilename(file.Name()))
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestCheckBackupShouldBeSkippedReturnTrueIfFileExists(t *testing.T) {
	b := createBackuper()
	b.OverwriteBackup = false
	file, err := ioutil.TempFile("", "backupfile_exists")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), []byte("backupcontent"), 0755)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	result := b.checkBackupShouldBeSkipped(model.BackupFilename("/tmp"))
	if err := AssertThat(result, Is(true)); err != nil {
		t.Fatal(err)
	}
}
