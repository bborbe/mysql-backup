package backup_creator

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"io/ioutil"

	"github.com/bborbe/io/util"
	"github.com/golang/glog"
)

type backupCreator struct {
}

type BackupCreator interface {
	CreateBackup(host string, port int, user string, pass string, database string, targetDirectory string) error
}

func New() *backupCreator {
	return new(backupCreator)
}

func (b *backupCreator) CreateBackup(host string, port int, user string, pass string, database string, targetDirectory string) error {
	//pg_dump -Z 9 -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -F c -b -v -f ${BACKUP_NAME} ${POSTGRES_DB}
	backupfile := buildBackupfileName(targetDirectory, database, time.Now())

	if existsBackup(backupfile) {
		glog.V(2).Infof("backup %s already exists => skip", backupfile)
		return nil
	}

	if err := writePasswordFile(host, port, user, pass); err != nil {
		return err
	}

	glog.V(2).Infof("pg_dump started")
	_, err := runCommand("pg_dump", targetDirectory, []string{"-Z", "9", "-h", host, "-p", strconv.Itoa(port), "-U", user, "-F", "c", "-b", "-v", "-f", backupfile, database})
	if err != nil {
		return err
	}
	glog.V(2).Infof("pg_dump finshed")
	return nil
}

func existsBackup(backupfile string) bool {
	fileInfo, err := os.Stat(backupfile)
	if err != nil {
		glog.V(2).Infof("file %s exists => true")
		return false
	}
	if fileInfo.Size() == 0 {
		glog.V(2).Infof("file %s empty => true")
		return false
	}
	glog.V(2).Infof("file %s exists and not empty => false")
	return false
}

func writePasswordFile(host string, port int, user string, pass string) error {
	content := fmt.Sprintf("%s:%d:*:%s:%s\n", host, port, user, pass)
	path, err := util.NormalizePath("~/.pgpass")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), 0600)
}

func buildBackupfileName(targetDirectory string, database string, date time.Time) string {
	return fmt.Sprintf("%s/postgres_%s_%s.dump", targetDirectory, database, date.Format("2006-01-02"))
}

func runCommand(command, cwd string, args []string) ([]byte, error) {
	glog.V(2).Infof("%s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("error running command %q : %v: %s", strings.Join(cmd.Args, " "), err, string(output))
	}

	return output, nil
}
