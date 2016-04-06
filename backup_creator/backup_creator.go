package backup_creator

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

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
	if _, err := os.Stat(backupfile); os.IsExist(err) {
		logger.Debugf("skip backup ")
		return nil
	}
	logger.Debugf("pg_dump started")
	_, err := runCommand("pg_dump", targetDirectory, []string{"-Z", "9", "-h", host, "-p", strconv.Itoa(port), "-U", user, "-F", "c", "-b", "-v", "-f", backupfile, database})
	if err != nil {
		return err
	}
	logger.Debugf("pg_dump finshed")
	return nil
}

func buildBackupfileName(targetDirectory string, database string, date time.Time) string {
	return fmt.Sprintf("%s/postgres_%s_%s.dump", targetDirectory, database, date.Format("2006-01-02"))
}

func runCommand(command, cwd string, args []string) ([]byte, error) {
	logger.Debugf("%s %s", command, strings.Join(args, " "))
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
