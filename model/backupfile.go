package model

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
)

type BackupFilename string

func BuildBackupfileName(targetDirectory TargetDirectory, database PostgresqlDatabase, date time.Time) BackupFilename {
	return BackupFilename(fmt.Sprintf("%s/postgres_%s_%s.dump", targetDirectory, database, date.Format("2006-01-02")))
}

func (b BackupFilename) Delete() error {
	return os.Remove(b.String())
}

func (b BackupFilename) String() string {
	return string(b)
}

func (b BackupFilename) Exists() bool {
	fileInfo, err := os.Stat(b.String())
	if err != nil {
		glog.V(2).Infof("file %v exists => false", b)
		return false
	}
	if fileInfo.Size() == 0 {
		glog.V(2).Infof("file %v empty => false", b)
		return false
	}
	glog.V(2).Infof("file %v exists and not empty => true", b)
	return true
}
