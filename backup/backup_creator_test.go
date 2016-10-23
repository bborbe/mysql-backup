package backup

import (
	"github.com/golang/glog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}
