package backup

import (
	. "github.com/bborbe/assert"
	"github.com/golang/glog"
	"os"
	"testing"
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
