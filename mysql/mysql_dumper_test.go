package mysql

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

func createDumper() *Dumper {
	return NewDumper(
		"database",
		"localhost",
		3306,
		"user",
		"password",
		"name",
		"/backup",
	)
}

func TestNew(t *testing.T) {
	mysqlDumper := createDumper()
	if err := AssertThat(mysqlDumper.Database.String(), Is("database")); err != nil {
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
