package backup

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

// Create backup
func Create(
	host string,
	port int,
	user string,
	pass string,
	database string,
	targetDirectory string,
) error {
	//pg_dump -Z 9 -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -F c -b -v -f ${BACKUP_NAME} ${POSTGRES_DB}
	backupfile := buildBackupfileName(targetDirectory, database, time.Now())

	if existsBackup(backupfile) {
		glog.V(1).Infof("backup %s already exists => skip", backupfile)
		return nil
	}

	if err := writePasswordFile(host, port, user, pass); err != nil {
		return err
	}

	glog.V(1).Infof("pg_dump started")
	if err := runCommand("pg_dump", targetDirectory, "-Z", "9", "-h", host, "-p", strconv.Itoa(port), "-U", user, "-F", "c", "-b", "-v", "-f", backupfile, database); err != nil {
		return err
	}
	glog.V(1).Infof("pg_dump finshed")
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

func runCommand(command, cwd string, args ...string) error {
	debug := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	glog.V(2).Infof("execute %s", debug)
	cmd := exec.Command(command, args...)
	if glog.V(4) {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if cwd != "" {
		cmd.Dir = cwd
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	glog.V(2).Infof("%s started", debug)
	if err := cmd.Wait(); err != nil {
		glog.Warningf("%s failed: %v", debug, err)
		return fmt.Errorf("%s failed: %v", debug, err)
	}
	glog.V(2).Infof("%s finished", command)
	return nil
}
