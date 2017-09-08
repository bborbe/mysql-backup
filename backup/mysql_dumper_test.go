package backup

import (
	"os"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func createDumper() *Backuper {
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
func createDumperAll() *Backuper {
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
	mysqlDumper := createDumper()
	if err := AssertThat(mysqlDumper.Database.String(), Is("database")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.AllDatabases, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Host.String(), Is("localhost")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Port.Int(), Is(3306)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.User.String(), Is("user")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Password.String(), Is("password")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Name.String(), Is("name")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.TargetDirectory.String(), Is("/backup")); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDumperAll(t *testing.T) {
	mysqlDumper := createDumperAll()
	if err := AssertThat(mysqlDumper.Database.String(), Is("")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.AllDatabases, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Host.String(), Is("localhost")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Port.Int(), Is(3306)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.User.String(), Is("user")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Password.String(), Is("password")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.Name.String(), Is("name")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(mysqlDumper.TargetDirectory.String(), Is("/backup")); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccess(t *testing.T) {
	mysqlDumper := createDumper()
	if err := AssertThat(mysqlDumper.Validate(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccessAll(t *testing.T) {
	mysqlDumper := createDumperAll()
	if err := AssertThat(mysqlDumper.Validate(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateReturnErrorIfDatabaseIsEmptyAndAllDatabasesActiv(t *testing.T) {
	mysqlDumper := createDumperAll()
	mysqlDumper.Database = "mydb"
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfDatabaseIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.Database = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfHostIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.Host = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfPortIsZero(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.Port = 0
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfUserIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.User = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfPasswordIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.Password = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfNameIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.Name = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateFailsIfTargetDirectoryIsEmpty(t *testing.T) {
	mysqlDumper := createDumper()
	mysqlDumper.TargetDirectory = ""
	if err := AssertThat(mysqlDumper.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
