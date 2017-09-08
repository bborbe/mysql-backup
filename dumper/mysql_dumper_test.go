package dumper

import (
	"io/ioutil"
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

func TestRunCommand(t *testing.T) {
	err := runCommand("ls", "/")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

var expectedContent = `
[mysqldump]
user=myuser
password=mypass
max_allowed_packet=1G
`

func TestWriteMyCnfFile(t *testing.T) {
	file, err := ioutil.TempFile("", "my.cnf")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = writeMyCnfFile(file.Name(), "myuser", "mypass")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadFile(file.Name())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(content), Is(expectedContent)); err != nil {
		t.Fatal(err)
	}
}
